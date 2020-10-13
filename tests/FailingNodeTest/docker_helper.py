import docker
from docker.models.resource import Model
import time


class DockerHelper:
    def __init__(self):
        self.client = docker.from_env()

    def get_name(self, name: str) -> Model:
        return self.client.containers.get(name)

    def restart_container(self, name: str) -> None:
        c = self.get_name(name).restart()

    def stop_container(self, name: str) -> None:
        c = self.get_name(name).stop()

    def start_container(self, name: str) -> None:
        self.get_name(name).start()

    def restart_container_timeout(self, name: str, timeout_s: int) -> None:
        self.stop_container(name)
        time.sleep(timeout_s)
        self.start_container(name)


if __name__ == '__main__':
    d = DockerHelper()
    print("Stopping Container nodeA-1")
    d.stop_container('nodeA-1')
    time.sleep(5)
    print("Restarting Container nodeA-1")
    d.start_container('nodeA-1')
