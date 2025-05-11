import bake

b = bake.Bakery()

GIT_ROOT = b.shell_strict("git rev-parse --show-toplevel")

DIRS = {
    "ROUTER" : f"{GIT_ROOT}/router",
    "SERVER" : f"{GIT_ROOT}/server",
    "SHARED" : f"{GIT_ROOT}/shared",
    "CLIENT" : f"{GIT_ROOT}/client",
    "DEPLOY" : f"{GIT_ROOT}/deploy",
    "BUILD" : f"{GIT_ROOT}/build",
}

GO = b.shell_strict("which go")

@b.target
def build() -> bool:
    """build all exec into build dir"""
    # ensure build dir
    b.shell_strict(f"mkdir -p {DIRS["BUILD"]}")

    dirs_execs = [
        (DIRS["ROUTER"], "router"),
        (DIRS["SERVER"], "server"),
        (DIRS["CLIENT"], "client"),
    ]

    def build_dir(route_exec):
        assert len(route_exec) == 2, f"invalid len(route_exec) : {len(route_exec)}"
        route = route_exec[0]
        exec = route_exec[1]
        b.shell_strict(f"cd {route} && {GO} build")
        b.shell_strict(f"mv {route}/{exec} {DIRS["BUILD"]}")

    # # build the router
    # b.shell_strict(f"cd {DIRS["ROUTER"]} && {GO} build")
    # b.shell_strict(f"mv {DIRS["ROUTER"]}/router {DIRS["BUILD"]}")

    # # build the server
    # b.shell_strict(f"cd {DIRS["SERVER"]} && {GO} build")
    # b.shell_strict(f"mv {DIRS["SERVER"]}/server {DIRS["BUILD"]}")

    # # build the client
    # b.shell_strict(f"cd {DIRS["CLIENT"]} && {GO} build")
    # b.shell_strict(f"mv {DIRS["CLIENT"]}/client {DIRS["BUILD"]}")

    _ = list(
        map(build_dir, dirs_execs)
    )

    return True

@b.target
def test() -> bool:
    """test all 3 workspaces"""

    b.shell_pass(f"go test {DIRS["ROUTER"]}/... -v")
    b.shell_pass(f"go test {DIRS["SERVER"]}/... -v")
    b.shell_pass(f"go test {DIRS["CLIENT"]}/... -v")
    b.shell_pass(f"go test {DIRS["SHARED"]}/... -v")

    return True

@b.target
def fmt() -> bool:
    """format all go code in the codebase"""

    format = lambda route: b.shell_pass(
        f"cd {route} && go fmt ./... && cd -"
    )

    _ = list(
        map(format, DIRS.values())
    )

    return True

@b.target
def docker_compose_full() -> bool:
    """clean compile docker_compose.yaml, start all servers"""

    docker_compose_file = f"{DIRS["DEPLOY"]}/docker_compose.yaml"

    b.shell_pass(f"rm {docker_compose_file}")

    pkl_command = " ".join([
        "pkl eval -f yaml",
        f"{DIRS["DEPLOY"]}/infrastructure.pkl > {docker_compose_file}"
    ])
    b.shell_strict(pkl_command)

    b.shell_strict(f"docker compose -f {docker_compose_file} up --build")

    return True

@b.target
def docker_compose_small() -> bool:
    """clean compile docker_compose.yaml from ./deploy/small_test.pkl, start all servers"""

    docker_compose_file = f"{DIRS["DEPLOY"]}/docker_compose.yaml"

    b.shell_pass(f"rm {docker_compose_file}")

    pkl_command = " ".join([
        "pkl eval -f yaml",
        f"{DIRS["DEPLOY"]}/small_test.pkl > {docker_compose_file}"
    ])
    b.shell_strict(pkl_command)

    b.shell_strict(f"docker compose -f {docker_compose_file} up --build")

    return True

if __name__ == "__main__":
    b.compile()
