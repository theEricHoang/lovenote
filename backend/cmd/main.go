package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/theEricHoang/lovenote/backend/internal/api"
)

func main() {
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

	r := api.RegisterRoutes()
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
