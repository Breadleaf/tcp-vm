import bake

b = bake.Bakery()

GIT_ROOT = b.shell_strict("git rev-parse --show-toplevel")

ROUTER_DIR = f"{GIT_ROOT}/router"
SERVER_DIR = f"{GIT_ROOT}/server"
SHARED_DIR = f"{GIT_ROOT}/shared"
CLIENT_DIR = f"{GIT_ROOT}/client"
DEPLOY_DIR= f"{GIT_ROOT}/deploy"
BUILD_DIR = f"{GIT_ROOT}/build"

GO = b.shell_strict("which go")

@b.target
def build() -> bool:
    """build all exec into build dir"""

    # ensure build dir
    b.shell_strict(f"mkdir -p {BUILD_DIR}")

    # build the router
    b.shell_strict(f"cd {ROUTER_DIR} && {GO} build")
    b.shell_strict(f"mv {ROUTER_DIR}/router {ROUTER_DIR}")

    # build the server
    b.shell_strict(f"cd {SERVER_DIR} && {GO} build")
    b.shell_strict(f"mv {SERVER_DIR}/server {BUILD_DIR}")

    # build the client
    b.shell_strict(f"cd {CLIENT_DIR} && {GO} build")
    b.shell_strict(f"mv {CLIENT_DIR}/client {BUILD_DIR}")

    return True

@b.target
def test() -> bool:
    """test all 3 workspaces"""

    b.shell_pass(f"go test {ROUTER_DIR}/... -v")
    b.shell_pass(f"go test {SERVER_DIR}/... -v")
    b.shell_pass(f"go test {CLIENT_DIR}/... -v")
    b.shell_pass(f"go test {SHARED_DIR}/... -v")

    return True

@b.target
def fmt() -> bool:
    """format all go code in the codebase"""

    b.shell_pass(f"gofmt -w {ROUTER_DIR}")
    b.shell_pass(f"gofmt -w {SERVER_DIR}")
    b.shell_pass(f"gofmt -w {CLIENT_DIR}")
    b.shell_pass(f"gofmt -w {SHARED_DIR}")

    return True

@b.target
def docker_compose_full() -> bool:
    """clean compile docker_compose.yaml, start all servers"""

    docker_compose_file = f"{DEPLOY_DIR}/docker_compose.yaml"

    b.shell_pass(f"rm {docker_compose_file}")

    pkl_command = " ".join([
        "pkl eval -f yaml",
        f"{DEPLOY_DIR}/infrastructure.pkl > {docker_compose_file}"
    ])
    b.shell_strict(pkl_command)

    b.shell_strict(f"docker compose -f {docker_compose_file} up --build")

    return True

@b.target
def docker_compose_small() -> bool:
    """clean compile docker_compose.yaml from ./deploy/small_test.pkl, start all servers"""

    docker_compose_file = f"{DEPLOY_DIR}/docker_compose.yaml"

    b.shell_pass(f"rm {docker_compose_file}")

    pkl_command = " ".join([
        "pkl eval -f yaml",
        f"{DEPLOY_DIR}/small_test.pkl > {docker_compose_file}"
    ])
    b.shell_strict(pkl_command)

    b.shell_strict(f"docker compose -f {docker_compose_file} up --build")

    return True

if __name__ == "__main__":
    b.compile()
