import time

from docker_helper import DockerHelper
from fred_client import FredClient
import sys

if __name__ == '__main__':
    dh = DockerHelper()
    nodeB = FredClient('172.26.2.1:9001')
    nodeB_name = 'nodeB'

    nodeC = FredClient('172.26.3.1:9001')
    nodeC_name = 'nodeC'
    total_errors = 0
    # Expected Status: Everything is running (make fred in 3NodeTest), nodeB is under the above IP
    print('------------------------------')
    print('Testing send/receive between nodes')

    nodeC.create_keygroup("basetest", True, 0)
    nodeC.add_replica("basetest", "nodeB", 0)
    nodeC.update("basetest", "key", "value")
    time.sleep(2)

    print('NodeB reading updated value')
    resp = nodeB.read("basetest", "key").data
    print(f'NodeB Read returned {resp}')

    if resp != "value":
        print(f'NodeC ERROR got resp {resp} instead of "value"')
        total_errors += 1

    print('------------------------------')
    print('Testing putting into node, restarting it and getting from same node')

    nodeC.create_keygroup('restarttest', True, 0)
    nodeC.update('restarttest', 'key', 'value')
    dh.restart_container_timeout(nodeC_name, 5)
    time.sleep(5)
    resp = nodeC.read('restarttest', 'key').data
    if resp != 'value':
        print(f'NodeC ERROR got resp {resp} instead of "value"')
        total_errors += 1

    print('------------------------------')
    print('Test: Create a replicated keygroup, kill nodeB, update the KG with nodeA, restart nodeB, test if nodeB has '
          'the data')

    nodeC.create_keygroup('disttest', True, 0)
    nodeC.add_replica('disttest', 'nodeB', 0)
    nodeC.update('disttest', 'key', 'value')
    dh.stop_container(nodeB_name)
    nodeC.update('disttest', 'key', 'newvalue')
    time.sleep(5)
    dh.start_container(nodeB_name)
    time.sleep(5)
    resp = nodeB.read('disttest', 'key').data
    if resp != 'newvalue':
        print(f'NodeB ERROR got resp {resp} instead of "newvalue"')
        total_errors += 1
        if resp == "value":
            print("...this means that nodeB has not received the updates while it was gone.")

    print("Done!")
    if total_errors > 0:
        print(f'Total errors: {total_errors}')
        sys.exit(total_errors)
