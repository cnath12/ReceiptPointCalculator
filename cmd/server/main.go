package main

import (
   "log"
   "net/http"
   "github.com/go-chi/chi/v5"
   "github.com/go-chi/chi/v5/middleware"
   "ReceiptPointCalculator/internal/api/handler"
   "ReceiptPointCalculator/internal/domain/service"
   "ReceiptPointCalculator/internal/storage/memory"
   "ReceiptPointCalculator/internal/validator"
   customMiddleware "ReceiptPointCalculator/internal/api/middleware"
)

func main() {
   // Initialize dependencies
   repo := memory.NewReceiptRepository()
   v := validator.NewReceiptValidator()
   svc := service.NewReceiptService(repo)
   handler := handler.NewReceiptHandler(svc, v)

   // Setup router
   r := chi.NewRouter()
   
   // Middleware
   r.Use(middleware.Logger)
   r.Use(middleware.Recoverer)
   r.Use(middleware.RequestID)
   r.Use(middleware.RealIP)
   r.Use(customMiddleware.ValidateRequest(v))

   // Routes
   r.Post("/receipts/process", handler.ProcessReceipt)
   r.Get("/receipts/{id}/points", handler.GetPoints)

   log.Println("Server starting on :8080")
   if err := http.ListenAndServe(":8080", r); err != nil {
       log.Fatal(err)
   }
}