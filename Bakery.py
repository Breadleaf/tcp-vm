import bake

b = bake.Bakery()

GIT_ROOT = b.shell_strict("git rev-parse --show-toplevel")

SERVER_DIR = f"{GIT_ROOT}/server"
CLIENT_DIR = f"{GIT_ROOT}/client"
SHARED_DIR = f"{GIT_ROOT}/shared"
BUILD_DIR = f"{GIT_ROOT}/build"

GO = b.shell_strict("which go")

@b.target
def build() -> bool:
    """build all exec into build dir"""

    # ensure build dir
    b.shell_strict(f"mkdir -p {BUILD_DIR}")

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

    b.shell_pass(f"go test {SERVER_DIR}/... -v")
    b.shell_pass(f"go test {CLIENT_DIR}/... -v")
    b.shell_pass(f"go test {SHARED_DIR}/... -v")

    return True

if __name__ == "__main__":
    b.compile()
