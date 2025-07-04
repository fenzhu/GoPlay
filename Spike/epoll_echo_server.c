#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <sys/epoll.h>
#include <errno.h>

#define MAX_EVENTS 10
#define PORT 9001

int main() {
    int server_fd, new_socket, epoll_fd;
    struct sockaddr_in address;
    int opt = 1;
    int addrlen = sizeof(address);
    char buffer[1024] = {0};

    // Create server socket
    if ((server_fd = socket(AF_INET, SOCK_STREAM, 0)) == 0) {
        perror("socket failed");
        exit(EXIT_FAILURE);
    }

    // Forcefully attaching socket to the port 8080
    if (setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR | SO_REUSEPORT, &opt, sizeof(opt))) {
        perror("setsockopt");
        exit(EXIT_FAILURE);
    }
    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(PORT);

    // Bind the socket to the network address and port
    if (bind(server_fd, (struct sockaddr *)&address, sizeof(address)) < 0) {
        perror("bind failed");
        exit(EXIT_FAILURE);
    }

    // Start listening for incoming connections
    if (listen(server_fd, 3) < 0) {
        perror("listen");
        exit(EXIT_FAILURE);
    }

    // Create epoll instance
    if ((epoll_fd = epoll_create1(0)) < 0) {
        perror("epoll_create1 failed");
        exit(EXIT_FAILURE);
    }

    struct epoll_event event, events[MAX_EVENTS];
    event.events = EPOLLIN; // Watch for incoming data
    event.data.fd = server_fd;

    // Add server socket to epoll interest list
    if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, server_fd, &event) < 0) {
        perror("epoll_ctl failed");
        exit(EXIT_FAILURE);
    }

    printf("Server listening on port %d\n", PORT);

    while(1) {
        int num_events = epoll_wait(epoll_fd, events, MAX_EVENTS, -1);
        if (num_events < 0) {
            perror("epoll_wait failed");
            exit(EXIT_FAILURE);
        }

        for (int i = 0; i < num_events; i++) {
            if (events[i].data.fd == server_fd) {
                // Event on the listening socket: accept new connection
                if ((new_socket = accept(server_fd, (struct sockaddr *)&address, (socklen_t*)&addrlen)) < 0) {
                    perror("accept failed");
                    continue;
                }
                printf("Accepted new connection on fd %d\n", new_socket);

                event.events = EPOLLIN | EPOLLET; // Use Edge-Triggered for client sockets
                event.data.fd = new_socket;
                if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, new_socket, &event) < 0) {
                    perror("epoll_ctl add client failed");
                    exit(EXIT_FAILURE);
                }
            } else {
                // Event on a client socket: read data
                int client_fd = events[i].data.fd;
                ssize_t count = read(client_fd, buffer, sizeof(buffer) - 1);

                if (count <= 0) {
                    // 0 means client closed connection, < 0 is an error
                    if (count == 0) {
                        printf("Client on fd %d disconnected\n", client_fd);
                    } else {
                        perror("read error");
                    }
                    epoll_ctl(epoll_fd, EPOLL_CTL_DEL, client_fd, NULL);
                    close(client_fd);
                } else {
                    // Echo data back to client
                    buffer[count] = '\0';
                    printf("Read %zd bytes from fd %d: %s", count, client_fd, buffer);
                    write(client_fd, buffer, count);
                }
            }
        }
    }

    close(server_fd);
    return 0;
}
