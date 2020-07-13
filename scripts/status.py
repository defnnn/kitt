import json
import os
import requests

# curl https://vault.$KITT_DOMAIN/v1/sys/seal-status
def status():
    status = {}
    url = "vault.{}".format(os.getenv("KITT_DOMAIN"))

    try:
        req = requests.get("https://{}/v1/sys/seal-status".format(url))
        req.raise_for_status()
        status = req.json()
    except requests.exceptions.RequestException as e:
        raise SystemExit(e)

    return status
