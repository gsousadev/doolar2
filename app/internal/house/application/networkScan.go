package application

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

func GetLocalIP() (string, *net.IPNet, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", nil, err
	}

	for _, iface := range ifaces {
		// Pula interfaces desligadas ou loopback
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			var ipnet *net.IPNet

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				ipnet = v
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // ignora IPv6
			}

			return ip.String(), ipnet, nil
		}
	}

	return "", nil, fmt.Errorf("não foi possível encontrar um IP local válido")
}

// ScanNetwork escaneia a rede local e retorna IP + MAC + nome
func ScanNetwork() error {
	localIP, ipnet, err := GetLocalIP()
	if err != nil {
		return err
	}

	fmt.Println("IP Local:", localIP)
	fmt.Println("Sub-rede:", ipnet.String())

	// Ping em toda a sub-rede para popular a tabela ARP local
	for ip := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		go func(ip string) {
			exec.Command("ping", "-c", "1", "-W", "1", ip).Run()
		}(ip.String())
	}
	time.Sleep(2 * time.Second) // Espera ping completar

	// Pega a tabela ARP
	out, err := exec.Command("arp", "-a").Output()
	if err != nil {
		return err
	}

	entries := strings.Split(string(out), "\n")
	for _, entry := range entries {
		if strings.TrimSpace(entry) == "" || strings.Contains(entry, "<incomplete>") {
			continue
		}

		fmt.Println("Dispositivo:", entry)
		// Exemplo de saída: ? (192.168.15.1) at aa:bb:cc:dd:ee:ff [ether] on en0
	}

	return nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func RunScanLoop() error {
	for {
		err := ScanNetwork()
		if err != nil {
			return err
		}
		fmt.Println("Varredura concluída. Aguardando 30 segundos para nova varredura...")
		time.Sleep(30 * time.Second)
	}
}
