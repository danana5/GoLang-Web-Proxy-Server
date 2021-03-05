import socket
import sys
from termcolor import colored, cprint
import _thread
import traceback
import ssl


def Proxy():
    global listen_port
    global buffer_size
    global max_conn

    listen_port = 8080
    max_conn = 10000

    try:
        print("[#] Initializing Socket.")
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        s.bind(("", listen_port))
        s.listen(max_conn)
        print("[#] Done.")
        print("[#] Socket has been binded successfully...")
        print("[#] Server working on: [{}]".format(listen_port))

    except Exception as e:
        print(e)
        sys.exit(2)


if __name__ == "__main__":
    Proxy()
