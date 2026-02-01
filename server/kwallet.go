package main

import (
	"fmt"
	"log"

	"github.com/godbus/dbus/v5"
)

// GetSecret retrieves a secret from KWallet via D-Bus
func GetSecret(key string) (string, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return "", fmt.Errorf("failed to connect to session bus: %v", err)
	}

	obj := conn.Object("org.kde.kwalletd5", "/modules/kwalletd5")

	var handle int32
	// open(wallet, wId, appid)
	// wallet="kdewallet", wId=0, appid="Dashboards"
	// Signature: (sxs)i  -> string, int64, string -> int32
	err = obj.Call("org.kde.KWallet.open", 0, "kdewallet", int64(0), "Dashboards").Store(&handle)
	if err != nil {
		return "", fmt.Errorf("failed to open wallet: %v", err)
	}

	var password string
	// readPassword(handle, folder, key, appid)
	// Signature: (isss)s -> int32, string, string, string -> string
	err = obj.Call("org.kde.KWallet.readPassword", 0, handle, "Dashboards", key, "Dashboards").Store(&password)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %v", err)
	}

	if password == "" {
		log.Printf("Warning: Secret for key '%s' is empty", key)
	}

	return password, nil
}
