package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
)

const grpcPort = "50051"

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func newInventoryService() *inventoryService {
	s := &inventoryService{parts: make(map[string]*inventoryV1.Part)}
	return s
}

func (s *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	log.Printf("GetPart: part with UUID %s found", req.GetUuid())
	return &inventoryV1.GetPartResponse{Part: part}, nil
}

func (s *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filter := req.GetFilter()

	parts := make([]*inventoryV1.Part, 0, len(s.parts))
	for _, part := range s.parts {
		if matches(part, filter) {
			parts = append(parts, part)
		}
	}

	// map перебирается в случайном порядке, поэтому сортируем для стабильного ответа.
	sort.Slice(parts, func(i, j int) bool { return parts[i].GetName() < parts[j].GetName() })

	return &inventoryV1.ListPartsResponse{Parts: parts}, nil
}

// matches: ИЛИ внутри одного поля фильтра, И между разными полями.
// Пустое поле фильтра не ограничивает результат.
func matches(part *inventoryV1.Part, f *inventoryV1.PartsFilter) bool {
	if f == nil {
		return true
	}

	if len(f.GetUuids()) > 0 && !contains(f.GetUuids(), part.GetUuid()) {
		return false
	}
	if len(f.GetNames()) > 0 && !contains(f.GetNames(), part.GetName()) {
		return false
	}
	if len(f.GetCategories()) > 0 && !containsCategory(f.GetCategories(), part.GetCategory()) {
		return false
	}
	if len(f.GetManufacturerCountries()) > 0 && !contains(f.GetManufacturerCountries(), part.GetManufacturer().GetCountry()) {
		return false
	}
	if len(f.GetTags()) > 0 && !intersects(f.GetTags(), part.GetTags()) {
		return false
	}

	return true
}

func contains(list []string, v string) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}

func containsCategory(list []inventoryV1.Category, v inventoryV1.Category) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}

// intersects: есть ли хотя бы один общий тег.
func intersects(a, b []string) bool {
	for _, x := range a {
		if contains(b, x) {
			return true
		}
	}
	return false
}

func main() {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	inventoryV1.RegisterInventoryServiceServer(server, newInventoryService())
	reflection.Register(server)

	go func() {
		log.Printf("🚀 gRPC InventoryService запущен на порту %s", grpcPort)
		if err := server.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы gRPC-сервера...")
	server.GracefulStop()
	log.Println("✅ Сервер остановлен")
}
