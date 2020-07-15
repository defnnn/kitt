import json
import os
import requests
import sys
sys.path.append('scripts/')
from pyinfra import host, local
from pyinfra.operations import python, server
from status import status

SHELL = '/usr/bin/env bash'

# curl --request PUT --data '{"secret_shares": 5, "secret_threshold": 3}' https://vault.$KITT_DOMAIN/v1/sys/init
def initialize():
    mStatus = status()
    init = {}

    url = "vault.{}".format(os.getenv("KITT_DOMAIN"))
    if "initialized" in mStatus and not mStatus["initialized"]:
        print('Initialize Moria')
        try:
            req = requests.put("https://{}/v1/sys/init".format(url), data = '{"secret_shares": 5, "secret_threshold": 3}')
            req.raise_for_status()
            init = req.json()
        except requests.exceptions.RequestException as e:
            raise SystemExit(e)

    return init

if host.fact.which('pass'):

    init = initialize()
    if init:
        print('Add Moria Keys to Pass')
        for index in init:
            if isinstance(init[index], list):
                if len(init[index]) > 0:
                    for objIndex, objItem in enumerate(init[index]):
                        server.shell({'Import moria key {0}_{1}'.format(index, objIndex + 1)}, 'echo {0} | pass insert -e moria/{1}_{2}'.format(objItem, index, objIndex + 1))
            else:
                server.shell({'Import moria key {}'.format(index)}, 'echo {0} | pass insert -e moria/{1}'.format(init[index], index))
        server.shell({'Pass git push'}, 'pass git push')

else:
    python.raise_exception(OSError, 'please ensure you have password-store installed and configured (passwordstore.org)')
