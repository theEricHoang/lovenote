package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/theEricHoang/lovenote/backend/internal/api"
)

func main() {
	fmt.Println(`
**                                               **          
/**                                              /**          
/**  ******  **    **  *****  *******   ******  ******  ***** 
/** **////**/**   /** **///**//**///** **////**///**/  **///**
/**/**   /**//** /** /******* /**  /**/**   /**  /**  /*******
/**/**   /** //****  /**////  /**  /**/**   /**  /**  /**//// 
***//******   //**   //****** ***  /**//******   //** //******
///  //////     //     ////// ///   //  //////     //   //////`)

	fmt.Printf("\nStarting server...\n")

	r := api.RegisterRoutes()
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
