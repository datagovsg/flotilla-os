job "zeppelin" {
  constraint {
    attribute = "${node.class}"
    value     = "l-cloud-nomad-client"
  }

  datacenters = ["ap-southeast-1a","ap-southeast-1b","ap-southeast-1c"]
  region      = "ap-southeast-1"
  type        = "service"

  update {
    max_parallel = 1
    auto_revert = true
  }

  group "zeppelin" {
    count = 1

    vault {
      policies = ["locus_aws_creds_spark_nomad","locus_spark_nomad_app"]
    }

    task "zeppelin" {
      driver = "docker"

      config {
        image      = "registry.gahmen.tech/l/ldatabase-zeppelin:change-base"
        force_pull = false

        network_mode = "host"

        volumes = [
          "alloc/spark-defaults.conf:/tmp/spark-defaults-addons.conf",
          "secrets/hive-site.xml:/opt/spark/conf/hive-site.xml",

          "alloc/shiro.ini.template:/zeppelin/conf/shiro.ini.template",
        ]

        mounts = [
          {
            source   = "fs-87f150c6/"
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
cat /tmp/spark-defaults-addons.conf >> ${SPARK_HOME}/conf/spark-defaults.conf; \
chown ${SPARK_USER} /zeppelin; \
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
        MAPBOX_ACCESS_TOKEN = "pk.eyJ1Ijoic2xhaTExIiwiYSI6ImNqbm9mZW1xMDJncGQzcHE5bmticjlwMmEifQ.6g5hcv9mqdmtimMUUqoeDQ"
        TZ                  = "Asia/Singapore"

        OAUTH2_CONFIG_CLIENT_ID = ""
        OAUTH2_CONFIG_SECRET    = ""
      }

      template {
        data = <<EOF
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

[users]
# List of users with their password allowed to access Zeppelin.
# To use a different strategy (LDAP / Database / ...) check the shiro doc at http://shiro.apache.org/configuration.html#Configuration-INISections
# To enable admin user, uncomment the following line and set an appropriate password.
admin = password1, admin
user1 = password2, role1, role2
user2 = password3, role3
user3 = password4, role2

# Sample LDAP configuration, for user Authentication, currently tested for single Realm
[main]
### A sample for configuring Active Directory Realm
#activeDirectoryRealm = org.apache.zeppelin.realm.ActiveDirectoryGroupRealm
#activeDirectoryRealm.systemUsername = userNameA

#use either systemPassword or hadoopSecurityCredentialPath, more details in http://zeppelin.apache.org/docs/latest/security/shiroauthentication.html
#activeDirectoryRealm.systemPassword = passwordA
#activeDirectoryRealm.hadoopSecurityCredentialPath = jceks://file/user/zeppelin/zeppelin.jceks
#activeDirectoryRealm.searchBase = CN=Users,DC=SOME_GROUP,DC=COMPANY,DC=COM
#activeDirectoryRealm.url = ldap://ldap.test.com:389
#activeDirectoryRealm.groupRolesMap = "CN=admin,OU=groups,DC=SOME_GROUP,DC=COMPANY,DC=COM":"admin","CN=finance,OU=groups,DC=SOME_GROUP,DC=COMPANY,DC=COM":"finance","CN=hr,OU=groups,DC=SOME_GROUP,DC=COMPANY,DC=COM":"hr"
#activeDirectoryRealm.authorizationCachingEnabled = false

### A sample for configuring LDAP Directory Realm
#ldapRealm = org.apache.zeppelin.realm.LdapGroupRealm
## search base for ldap groups (only relevant for LdapGroupRealm):
#ldapRealm.contextFactory.environment[ldap.searchBase] = dc=COMPANY,dc=COM
#ldapRealm.contextFactory.url = ldap://ldap.test.com:389
#ldapRealm.userDnTemplate = uid={0},ou=Users,dc=COMPANY,dc=COM
#ldapRealm.contextFactory.authenticationMechanism = simple

### A sample PAM configuration
#pamRealm=org.apache.zeppelin.realm.PamRealm
#pamRealm.service=sshd

### A sample for configuring ZeppelinHub Realm
#zeppelinHubRealm = org.apache.zeppelin.realm.ZeppelinHubRealm
## Url of ZeppelinHub
#zeppelinHubRealm.zeppelinhubUrl = https://www.zeppelinhub.com
#securityManager.realms = $zeppelinHubRealm

## A same for configuring Knox SSO Realm
#knoxJwtRealm = org.apache.zeppelin.realm.jwt.KnoxJwtRealm
#knoxJwtRealm.providerUrl = https://domain.example.com/
#knoxJwtRealm.login = gateway/knoxsso/knoxauth/login.html
#knoxJwtRealm.logout = gateway/knoxssout/api/v1/webssout
#knoxJwtRealm.logoutAPI = true
#knoxJwtRealm.redirectParam = originalUrl
#knoxJwtRealm.cookieName = hadoop-jwt
#knoxJwtRealm.publicKeyPath = /etc/zeppelin/conf/knox-sso.pem
#
#knoxJwtRealm.groupPrincipalMapping = group.principal.mapping
#knoxJwtRealm.principalMapping = principal.mapping
#authc = org.apache.zeppelin.realm.jwt.KnoxAuthenticationFilter

sessionManager = org.apache.shiro.web.session.mgt.DefaultWebSessionManager

### If caching of user is required then uncomment below lines
#cacheManager = org.apache.shiro.cache.MemoryConstrainedCacheManager
#securityManager.cacheManager = $cacheManager

### Enables 'HttpOnly' flag in Zeppelin cookies
cookie = org.apache.shiro.web.servlet.SimpleCookie
cookie.name = JSESSIONID
cookie.httpOnly = true
### Uncomment the below line only when Zeppelin is running over HTTPS
#cookie.secure = true
sessionManager.sessionIdCookie = $cookie

securityManager.sessionManager = $sessionManager
# 86,400,000 milliseconds = 24 hour
securityManager.sessionManager.globalSessionTimeout = 86400000
shiro.loginUrl = /api/login

[roles]
role1 = *
role2 = *
role3 = *
admin = *

[urls]
# This section is used for url-based security. For details see the shiro.ini documentation.
#
# You can secure interpreter, configuration and credential information by urls.
# Comment or uncomment the below urls that you want to hide:
# anon means the access is anonymous.
# authc means form based auth Security.
#
# IMPORTANT: Order matters: URL path expressions are evaluated against an incoming request
# in the order they are defined and the FIRST MATCH WINS.
#
# To allow anonymous access to all but the stated urls,
# uncomment the line second last line (/** = anon) and comment the last line (/** = authc)
#
/api/version = anon
# Allow all authenticated users to restart interpreters on a notebook page.
# Comment out the following line if you would like to authorize only admin users to restart interpreters.
/api/interpreter/setting/restart/** = authc
/api/interpreter/** = authc, roles[admin]
/api/configurations/** = authc, roles[admin]
/api/credential/** = authc, roles[admin]
#/** = anon
/** = authc

EOF

        destination = "alloc/shiro.ini.template"
      }

      template {
        data = <<EOF
spark.hadoop.fs.s3.impl                                org.apache.hadoop.fs.s3a.S3AFileSystem
spark.hadoop.fs.s3a.impl                               org.apache.hadoop.fs.s3a.S3AFileSystem
spark.hadoop.fs.s3a.server-side-encryption-algorithm   SSE-KMS
spark.hadoop.fs.s3a.server-side-encryption.key         {{ key "terraform/spark_nomad/kms_alias_aws_s3_key_path" }}

# Need to be set on Spark master, workers and drivers
spark.ui.reverseProxy                              true
spark.ui.reverseProxyUrl                           https://sm.locus.rocks

spark.driver.bindAddress                           0.0.0.0
spark.driver.host                                  {{ with node }}{{ .Node.Address }}{{ end }}

spark.cores.max                                    {{ key "terraform/spark_nomad/cores_max_path" }}
spark.default.parallelism                          {{ key "terraform/spark_nomad/default_parallelism_path" }}

spark.sql.catalogImplementation                    {{ key "terraform/spark_nomad/sql_catalog_implementation" }}
spark.sql.warehouse.dir                            {{ key "terraform/spark_nomad/hive_metastore_s3_warehouse" }}

spark.port.maxRetries                              {{ key "terraform/spark_nomad/port_max_retries" }}

spark.shuffle.service.enabled                      {{ key "terraform/spark_nomad/shuffle_service_enabled" }}

spark.dynamicAllocation.enabled                    {{ key "terraform/spark_nomad/dynamic_allocation_enabled" }}
spark.dynamicAllocation.initialExecutors           {{ key "terraform/spark_nomad/dynamic_allocation_initial_executors" }}
spark.dynamicAllocation.minExecutors               {{ key "terraform/spark_nomad/dynamic_allocation_min_executors" }}
spark.dynamicAllocation.maxExecutors               {{ key "terraform/spark_nomad/dynamic_allocation_max_executors" }}
spark.dynamicAllocation.executorIdleTimeout        {{ key "terraform/spark_nomad/dynamic_allocation_executor_idle_timeout" }}
spark.dynamicAllocation.cachedExecutorIdleTimeout  {{ key "terraform/spark_nomad/dynamic_allocation_cached_executor_idle_timeout" }}

spark.executor.instances                           {{ key "terraform/spark_nomad/executor_instances" }}
spark.executor.cores                               {{ key "terraform/spark_nomad/executor_cores" }}
spark.executor.memory                              {{ key "terraform/spark_nomad/executor_memory" }}

EOF

        destination = "alloc/spark-defaults.conf"
      }

      template {
        data = <<EOF
<?xml version="1.0"?>
<configuration>

<property>
  <name>hive.metastore.pre.event.listeners</name>
  <value></value>
</property>
<property>
  <name>hive.metastore.warehouse.dir</name>
  <value>{{ key "terraform/spark_nomad/hive_metastore_s3_warehouse" }}</value>
</property>
<property>
  <name>hive.mv.files.thread</name>
  <value>15</value>
</property>
<property>
  <name>hive.warehouse.subdir.inherit.perms</name>
  <value>false</value>
</property>
<property>
  <name>javax.jdo.option.ConnectionURL</name>
  <value>{{ key "terraform/spark_nomad/hive_site_conf_db_connection_url" }}</value>
</property>
<property>
  <name>javax.jdo.option.ConnectionDriverName</name>
  <value>org.mariadb.jdbc.Driver</value>
</property>
{{ with secret "locus/app/spark_nomad" }}
<property>
  <name>javax.jdo.option.ConnectionUserName</name>
  <value>{{ .Data.hive_metastore_user }}</value>
</property>
<property>
  <name>javax.jdo.option.ConnectionPassword</name>
  <value>{{ .Data.hive_metastore_password }}</value>
</property>
{{ end }}

</configuration>

EOF

        destination = "secrets/hive-site.xml"
      }

      template {
        data = <<EOF
{{ with secret "locus_aws/creds/spark_nomad" }}
AWS_ACCESS_KEY_ID="{{ .Data.access_key }}"
AWS_SECRET_ACCESS_KEY="{{ .Data.secret_key }}"
AWS_DEFAULT_REGION="ap-southeast-1"
{{ end }}

PORT="{{ env "NOMAD_HOST_PORT_webui" }}"
NOMAD_TOKEN="{{ with secret "locus_nomad/creds/spark_nomad" }}{{ .Data.secret_id }}{{ end }}"
SPARK_MASTER="spark://{{ with service "spark-master" }}{{ with index . 0 }}{{ .Address }}:{{ .Port }}{{ end }}{{ end }}"
EOF
        destination = "secrets/env"
        env         = true
      }

      service {
        name = "${JOB}"
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
          "traefik.frontend.rule=Host:zeppelin.locus.rocks",
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
