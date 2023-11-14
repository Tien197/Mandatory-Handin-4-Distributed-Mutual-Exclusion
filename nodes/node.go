package main

import (
	"bufio"
	"context"
	"flag"
	"github.com/Tien197/Mandatory-Handin-4-Distributed-Mutual-Exclusion/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"strconv"
)

type Node struct {
	proto.UnimplementedNodeServer
	id               int
	port             int
	timestamp        int
	queue            map[int]int // store node id and node timestamp
	hasEntered       bool
	currentlyWaiting bool
	allNodes         map[int]proto.NodeClient // store individual node info
}

/*
go run nodes/node.go -port 8080 -id 1
*/
var (
	port = flag.Int("port", 0, "client port number")
	id   = flag.Int("id", 0, "client ID number")
)

func main() {
	flag.Parse()

	// Create a client
	node := &Node{
		id:               *id,
		port:             *port,
		timestamp:        1,
		queue:            make(map[int]int),
		hasEntered:       false,
		currentlyWaiting: false,
		allNodes:         make(map[int]proto.NodeClient),
	}

	// Initialize nodes with predetermined ports
	node.allNodes[1], _ = connectToNode(8080)
	node.allNodes[2], _ = connectToNode(8081)
	node.allNodes[3], _ = connectToNode(8082)

	// Starts the client
	go startNode(node)

	go waitForRequest(node)

	// keeps node running until exit/disconnect
	for {

	}
}

func startNode(node *Node) {
	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(node.port))

	if err != nil {
		log.Fatalf("Could not create the node %v", err)
	} else {
		log.Printf("Started node %d at port: %d at Lamport timestamp %d\n", node.id, node.port, node.timestamp)
	}

	proto.RegisterNodeServer(grpcServer, node)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func waitForRequest(node *Node) {

	// Wait for input in the node terminal
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		node.timestamp++
		log.Printf("Node requests to enter Critical Section at timestamp %d\n", node.timestamp)

		node.queue[node.id] = node.timestamp
		node.currentlyWaiting = true

		// current node asks for every node in the queue for permission to join

		// for every node in the queue (other than itself), it compares its timestamp with the other nodes

		// if current node's timestamp is lower (meaning that it was first to request) than the comparing node's timestamp
		// then it shall be given permission to enter critical section

		// set currentlyWaiting = false
		// remember to set hasEntered = true when entering

		// if not then it has to continue to wait

		for otherId, nodeConn := range node.allNodes {
			if otherId != node.id {
				node.timestamp++
				log.Printf("Node %d asks Node %d for permission to enter Critical Section at lamport timestamp %d \n", node.id, otherId, node.timestamp)

				_, _ = nodeConn.RequestToEnterSection(context.Background(), &proto.NodeMessage{
					Id:        int64(node.id),
					Timestamp: int64(node.timestamp),
				})

				/*if err != nil {
					log.Printf("%v", err)
				}*/
			}
		}

	}
}

// Add the request to the local queue, update local timestamp, and reply accordingly
func (s *Node) RequestToEnter(ctx context.Context, in *proto.NodeMessage) (*proto.NodeMessage, error) {

	return &proto.NodeMessage{}, nil
}

// Check if the request can be granted and reply accordingly
func (s *Node) EnterSection(ctx context.Context, in *proto.NodeMessage) (*proto.NodeMessage, error) {

	return &proto.NodeMessage{}, nil
}

// Remove node from critical section and grant access to the next in queue if any
func (s *Node) LeaveSection(ctx context.Context, in *proto.NodeMessage) (*proto.NodeMessage, error) {

	return &proto.NodeMessage{}, nil
}

func connectToNode(port int) (proto.NodeClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", port)
	}
	return proto.NewNodeClient(conn), nil
}
