//go:build windows

package main

import (
	"fmt"
	"log"
	"syscall"
	_ "unsafe"

	"golang.org/x/sys/windows"
)

func main() {
	fmt.Println("[+] Starting basic Go DLL loader...")

	dllPath := "calc_dll.dll"
	fmt.Printf("[+] Attempting to load DLL: %s\n", dllPath)

	dllHandle, err := windows.LoadLibrary(dllPath)
	if err != nil {
		log.Fatalf("[-] Failed to load DLL '%s': %v\n", dllPath, err)
	}

	defer func() {
		fmt.Println("[+] Attempting to free DLL handle...")
		err := windows.FreeLibrary(dllHandle)
		if err != nil {
			log.Printf("[!] Warning: Failed to free DLL handle: %v\n", err)
		} else {
			fmt.Println("[+] DLL handle freed successfully.")
		}
	}()

	fmt.Printf("[+] DLL loaded successfully. Handle: 0x%X\n", dllHandle)

	funcName := "LaunchCalc"
	fmt.Printf("[+] Attempting to get address of function: %s\n", funcName)

	funcAddr, err := windows.GetProcAddress(dllHandle, funcName)
	if err != nil {
		log.Fatalf("[-] Failed to find function '%s' in DLL: %v\n", funcName, err)
	}

	fmt.Printf("[+] Function '%s' found at address: 0x%X\n", funcName, funcAddr)

	fmt.Printf("[+] Calling function '%s'...\n", funcName)

	ret, _, callErr := syscall.SyscallN(funcAddr, 0, 0, 0, 0)

	if callErr != 0 {
		log.Fatalf("[-] Error occurred during syscall to '%s': %v\n", funcName, callErr)
	}

	if ret != 0 {
		fmt.Printf("[+] Function '%s' executed successfully (returned TRUE).\n", funcName)
	} else {
		fmt.Printf("[-] Function '%s' execution reported failure (returned FALSE).\n", funcName)
	}

	fmt.Println("[+] Loader finished.")
}
