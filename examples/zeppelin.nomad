job "zeppelin" {
  constraint {
    attribute = "$${node.class}"
    value     = "${node_class}"
  }

  datacenters = ${az}
  region      = "${region}"
  type        = "service"

  update {
    max_parallel = 1
    auto_revert = true
  }

  group "zeppelin" {
    count = 1

    vault {
      policies = ${vault_policies}
    }

    task "zeppelin" {
      driver = "docker"

      config {
        image      = "${docker_image}:${docker_tag}"
        force_pull = false

        network_mode = "host"

        volumes = [
          "alloc/spark-defaults.conf:/tmp/spark-defaults-addons.conf",
          "secrets/hive-site.xml:/opt/spark/conf/hive-site.xml",

          "alloc/shiro.ini.template:/zeppelin/conf/shiro.ini.template",
        ]

        mounts = [
          {
            source   = "${efs_id}/"
            target   = "/mnt/efs"
            readonly = false

            volume_options {
              driver_config {
                name = "efs"
              }
            }
          },
        ]

        privileged = true

        command = "sh"
        args = [
          "-c", <<EOF
set -euo pipefail; \
cat /tmp/spark-defaults-addons.conf >> $${SPARK_HOME}/conf/spark-defaults.conf; \
chown $${SPARK_USER} /zeppelin; \
mkdir -p /mnt/efs/zeppelin/notebook; \
mount --bind --make-rprivate /mnt/efs/zeppelin/notebook /zeppelin/notebook; \
mkdir -p /mnt/efs/locus; \
mount --bind --make-rprivate /mnt/efs/zeppelin/notebook /locus; \
exec bash /zeppelin/conf/run-zeppelin.sh; \
:
EOF
        ]

        dns_servers = ["169.254.1.1"]
      }

      env {
        NOMAD_ADDR = "https://nomad.locus.rocks"

        SPARK_USER          = "zeppelin"
        PUBLIC_SERVER       = "https://zep.locus.rocks"
        MAPBOX_ACCESS_TOKEN = "${mapbox_access_token}"
        TZ                  = "Asia/Singapore"

        OAUTH2_CONFIG_CLIENT_ID = ""
        OAUTH2_CONFIG_SECRET    = ""
      }

      template {
        data = <<EOF
${shiro_ini_template}
EOF

        destination = "alloc/shiro.ini.template"
      }

      template {
        data = <<EOF
${spark_defaults_conf}
EOF

        destination = "alloc/spark-defaults.conf"
      }

      template {
        data = <<EOF
${hive_site_conf}
EOF

        destination = "secrets/hive-site.xml"
      }

      template {
        data = <<EOF
{{ with secret "${spark_nomad_aws_creds_path}" }}
AWS_ACCESS_KEY_ID="{{ .Data.access_key }}"
AWS_SECRET_ACCESS_KEY="{{ .Data.secret_key }}"
AWS_DEFAULT_REGION="${region}"
{{ end }}

PORT="{{ env "NOMAD_HOST_PORT_webui" }}"
NOMAD_TOKEN="{{ with secret "${spark_nomad_nomad_path}" }}{{ .Data.secret_id }}{{ end }}"
SPARK_MASTER="spark://{{ with service "spark-master" }}{{ with index . 0 }}{{ .Address }}:{{ .Port }}{{ end }}{{ end }}"
EOF
        destination = "secrets/env"
        env         = true
      }

      service {
        name = "$${JOB}"
        port = "webui"

        check {
          type     = "http"
          path     = "/"
          interval = "10s"
          timeout  = "2s"

          check_restart {
            limit = 3
            grace = "60s"
          }
        }

        # Refer to https://docs.traefik.io/configuration/backends/consulcatalog/ for documentation
        # Be especially careful when setting the Content Security Policy.
        # See https://content-security-policy.com/
        tags = [
          "traefik.enable=true",
          "traefik.frontend.rule=Host:${zeppelin_domain}",
          "traefik.frontend.entryPoints=internal",
          "traefik.frontend.headers.SSLRedirect=true",
          "traefik.frontend.headers.SSLProxyHeaders=X-Forwarded-Proto:https",
          "traefik.frontend.headers.STSSeconds=315360000",
          "traefik.frontend.headers.frameDeny=true",
          "traefik.frontend.headers.browserXSSFilter=true",
          "traefik.frontend.headers.contentTypeNosniff=true",
          "traefik.frontend.headers.referrerPolicy=strict-origin",
          "traefik.frontend.headers.contentSecurityPolicy=default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' blob:; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob:; connect-src 'self'; font-src 'self' data:;",
        ]
      }

      resources {
        cpu    = 2000
        memory = 4096

        network {
          mbits = 100

          port "webui" {}
        }
      }
    }
  }
}
