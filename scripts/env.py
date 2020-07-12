import os
from pyinfra import host, local
from pyinfra.operations import python

path = os.path.split(os.path.dirname(os.path.realpath(__file__)))[0]

if host.fact.which('pass'):

    if os.path.exists("{}/.env".format(path)):
        try:
            os.remove("{}/.env".format(path))
        except:
            python.raise_exception(OSError, 'error deleting legacy kitt .env file')

    qass = local.shell('pass kitt | tail -n +2')
    for index in qass.splitlines():
        key = index.split()[1]
        value = local.shell('pass kitt/{}'.format(key))
        with open("{}/.env".format(path), "a") as env:
                env.write('{0}={1}\n'.format(key, value))

else:
    python.raise_exception(OSError, 'please ensure you have password-store installed and configured (passwordstore.org)')
