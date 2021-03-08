import socket
import sys
import _thread
import traceback
import ssl


def Proxy():
    global listen_port
    global buffer_size
    global max_conn

    listen_port = 8080
    max_conn = 10000
    buffer_size = 10000

    try:
        print("[#] Initializing Socket.")
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        s.bind(("127.0.0.1", listen_port))
        s.listen()
        print("[#] Done baby x")
        print("[#] Socket has been binded successfully...")
        print("[#] Server working on: [{}]".format(listen_port))

    except Exception as e:
        print(e)
        sys.exit(2)

    while 1 < 2:
        try:
            conn, addr = s.accept()
            data = conn. recv(buffer_size)
            _thread.start_new_thread(conn_string, (conn, data, addr))
        except KeyboardInterrupt:
            s.close()
            print("Goodbye My Lover <3 <3 <3 <3...")
            sys.exit(1)

    s.close()


def conn_string(conn, data, addr):
    try:
        print(addr)
        first_line = data.decode("utf-8").split("\n")[0]
        print(first_line)
        url = first_line.split(" ")[1]

        http_pos = url.find("://")
        if (http_pos == -1):
            temp = url
        else:
            temp = url[(http_pos + 3):]

        port_pos = temp.find("/")
        webserver_pos = temp.find("/")

        if(webserver_pos == -1):
            webserver_pos = len(temp)

        webserver = ""
        port = -1

        if (port_pos == -1 or webserver_pos < port_pos):
            port = 80
            webserver = temp[:webserver_pos]
        else:
            port = int(temp[(port_pos + 1):][:webserver_pos - port_pos - 1])
            webserver = temp[:port_pos]

        print("webserver= " + webserver)
        proxy_server(webserver, port, conn, data, addr)

    except Exception as e:
        print(e)
        traceback.print_exc()


def proxy_server(webserver, port, conn, data, addr):
    print("{}".format(webserver))
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((webserver, port))
        s.send(data)
        while 1:
            reply = s.recv(buffer_size)

            if len(reply) > 0:
                conn.sendall(reply)
                print("[*] Request sent: {} > {}".format(addr[0], webserver))
            else:
                break

        s.close()
        conn.close()

    except Exception as e:
        print(e)
        traceback.print_exc()
        s.close()
        conn.close()
        sys.exit(1)


if __name__ == "__main__":
    Proxy()
