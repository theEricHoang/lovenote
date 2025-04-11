package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	config "github.com/theEricHoang/lovenote/backend/internal"
	"github.com/theEricHoang/lovenote/backend/internal/api"
	"github.com/theEricHoang/lovenote/backend/internal/api/auth"
	"github.com/theEricHoang/lovenote/backend/internal/api/middleware"
	"github.com/theEricHoang/lovenote/backend/internal/api/users/dao"
	"github.com/theEricHoang/lovenote/backend/internal/api/users/handlers"
	"github.com/theEricHoang/lovenote/backend/internal/pkg/db"
	"github.com/theEricHoang/lovenote/backend/internal/pkg/imageservice"

	notedao "github.com/theEricHoang/lovenote/backend/internal/api/notes/dao"
	notehandlers "github.com/theEricHoang/lovenote/backend/internal/api/notes/handlers"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	config.LoadConfig() // load env vars

	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load aws config: %v", err)
	}
	s3Client := s3.NewFromConfig(awsCfg)
	presignClient := s3.NewPresignClient(s3Client)
	presigner := imageservice.NewPresigner(presignClient)

	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	authService := auth.NewAuthService(database)
	userDAO := dao.NewUserDAO(database)
	relationshipDAO := dao.NewRelationshipDAO(database)
	inviteDAO := dao.NewInviteDAO(database)
	noteDAO := notedao.NewNoteDAO(database)

	authMiddleware := middleware.NewAuthMiddleware(authService)
	permissionsMiddleware := middleware.NewPermissionsMiddleware(relationshipDAO, noteDAO)

	userHandler := handlers.NewUserHandler(userDAO, authService)
	relationshipHandler := handlers.NewRelationshipHandler(relationshipDAO)
	inviteHandler := handlers.NewInviteHandler(inviteDAO, relationshipDAO)
	noteHandler := notehandlers.NewNoteHandler(noteDAO, relationshipDAO)

	// shutdown signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// handle shutdown gracefully
	go func() {
		<-c // wait for shutdown signal to be received
		fmt.Println("\nShutting down gracefully...")
		database.Close()
		os.Exit(0)
	}()

	// ANSI escape codes for colors
	red := "\033[1;31m"
	reset := "\033[0m"

	fmt.Println(red + `
**                                               **          
/**                                              /**          
/**  ******  **    **  *****  *******   ******  ******  ***** 
/** **////**/**   /** **///**//**///** **////**///**/  **///**
/**/**   /**//** /** /******* /**  /**/**   /**  /**  /*******
/**/**   /** //****  /**////  /**  /**/**   /**  /**  /**//// 
***//******   //**   //****** ***  /**//******   //** //******
///  //////     //     ////// ///   //  //////     //   //////` + reset)

	fmt.Printf("\n\tStarting server, listening at port :8000...\n\n")

	r := api.RegisterRoutes(userHandler, relationshipHandler, inviteHandler, noteHandler, authMiddleware, permissionsMiddleware, presigner)
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
