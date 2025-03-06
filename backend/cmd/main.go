package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/theEricHoang/lovenote/backend/internal/api"
	"github.com/theEricHoang/lovenote/backend/internal/api/users/dao"
	"github.com/theEricHoang/lovenote/backend/internal/api/users/handlers"
	db "github.com/theEricHoang/lovenote/backend/internal/pkg"
)

func main() {
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	userDAO := dao.NewUserDAO(database)
	relationshipDAO := dao.NewRelationshipDAO(database)
	inviteDAO := dao.NewInviteDAO(database)

	userHandler := handlers.NewUserHandler(userDAO)
	relationshipHandler := handlers.NewRelationshipHandler(relationshipDAO)
	inviteHandler := handlers.NewInviteHandler(inviteDAO, relationshipDAO)

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

	r := api.RegisterRoutes(userHandler, relationshipHandler, inviteHandler)
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
