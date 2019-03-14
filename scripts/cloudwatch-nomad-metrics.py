#!/usr/bin/env python3

# adapted from https://github.com/trackit/cloudwatch-nomad-metrics

import datetime
import sched
import time
import json

from copy import deepcopy
from pprint import pprint

import nomad
import boto3

# config for this instance assumes that you have aws_credentials
client_events = boto3.client('events')

# config for this instance assumes that you have NOMAD_TOKEN and NOMAD_ADDR in envvar
# export NOMAD_TOKEN=`cat ~/.nomad-token`; export NOMAD_ADDR=https://nomad.locus.rocks
client_nomad = nomad.Nomad()

# Source = "nomad.script", DetailType = 'Nomad Job State Change'
# is the key to how cloudwatch re-routes the json objects to SQS
json_template = {
    'Source': 'nomad.script',
    'Resources': ['resource1'],
    'DetailType': 'Nomad Job State Change',
}

JOB_PREFIX = "flotilla"

def put_events(now):
    jobs = client_nomad.jobs.get_jobs()
    names = [j["Name"] for j in jobs if JOB_PREFIX in j["Name"]]
    entries = []
    for name in names:
        job = client_nomad.job.get_job(name)
        entry = deepcopy(json_template)
        entry["Detail"] = json.dumps(job)
        entries.append(entry)
    pprint(job)
    res = client_events.put_events(
        Entries=entries,
    )
    print("just delivered jobs to queue")
    enter_next(s, put_events)

def put_events2():
    jobs = client_nomad.jobs.get_jobs()
    names = [j["Name"] for j in jobs if JOB_PREFIX in j["Name"]]
    entries = []
    for name in names:
        job = client_nomad.job.get_job(name)
        entry = deepcopy(json_template)
        entry["Detail"] = json.dumps(job)
        entries.append(entry)
    if entries != []:
        pprint(job)
        res = client_events.put_events(
            Entries=entries,
        )
        print("just delivered jobs to queue")
    else:
        print("no jobs with {} prefix".format(JOB_PREFIX))
    time.sleep(10)

def enter_next(s, function):
    now = datetime.datetime.now()
    now_next = now.replace(
        second=0,
        microsecond=0,
    ) + datetime.timedelta(minutes=1)
    s.enterabs(
        time=now_next.timestamp(),
        priority=1,
        action=function,
        argument=(now_next,),
    )


if __name__ == '__main__':
    ## s = sched.scheduler(lambda: datetime.datetime.now().timestamp(), time.sleep)
    ## enter_next(s, put_events)
    while True:
        # s.run()
        put_events2()
