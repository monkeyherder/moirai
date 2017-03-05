package config_test

import (
	. "github.com/monkeyherder/moirai/config"

	"encoding/json"
	"fmt"
	"github.com/monkeyherder/moirai/checks"
	"github.com/monkeyherder/moirai/checks/network"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("ChecksConfig", func() {

	Context("Given a valid checksconfig", func() {
		var checksdConfig *ChecksdConfig
		var checksdConfigJson []byte

		BeforeEach(func() {
			checksdConfig = &ChecksdConfig{
				ChecksPollTime: 1 * time.Second,
			}
			var err error
			checksdConfigJson, err = json.Marshal(checksdConfig)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("with network icmp check config", func() {
			BeforeEach(func() {
				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "icmp",
						CheckProperties: map[string]interface{}{
							"Address": "icmp_address",
							"Timeout": 1 * time.Second,
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)

				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal valid icmpcheck", func() {
				err := json.Unmarshal(checksdConfigJson, checksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(checksdConfig.Checks).To(HaveLen(1))
				Expect(checksdConfig.Checks[0]).To(Equal(&network.IcmpCheck{
					Address: "icmp_address",
					Timeout: 1 * time.Second,
				}))
			})
		})

		Context("with network socket check config", func() {
			BeforeEach(func() {
				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "unix_socket",
						CheckProperties: map[string]interface{}{
							"Timeout":    1 * time.Second,
							"SocketFile": "unix-socket-file",
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal valid socket check", func() {
				err := json.Unmarshal(checksdConfigJson, checksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(checksdConfig.Checks).To(HaveLen(1))
				Expect(checksdConfig.Checks[0]).To(Equal(&network.UnixSocketCheck{
					Timeout:    1 * time.Second,
					SocketFile: "unix-socket-file",
				}))
			})
		})

		Context("with network tcp check config", func() {
			BeforeEach(func() {
				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "tcp",
						CheckProperties: map[string]interface{}{
							"Timeout": 1 * time.Second,
							"Port":    1234,
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal valid tcp check", func() {
				err := json.Unmarshal(checksdConfigJson, checksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(checksdConfig.Checks).To(HaveLen(1))
				Expect(checksdConfig.Checks[0]).To(Equal(&network.TcpCheck{
					Timeout: 1 * time.Second,
					Port:    1234,
				}))
			})
		})

		Context("with network udp check config", func() {
			BeforeEach(func() {
				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "udp",
						CheckProperties: map[string]interface{}{
							"Timeout": 1 * time.Second,
							"Port":    1234,
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal valid udp check", func() {
				err := json.Unmarshal(checksdConfigJson, checksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(checksdConfig.Checks).To(HaveLen(1))
				Expect(checksdConfig.Checks[0]).To(Equal(&network.UdpCheck{
					Timeout: 1 * time.Second,
					Port:    1234,
				}))
			})
		})

		Context("with file check config", func() {
			BeforeEach(func() {

				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "file",
						CheckProperties: map[string]interface{}{
							"Name":      "name",
							"Path":      "path",
							"IfChanged": "ifchanged",
							"Group":     "group",
							"DependsOn": "dependson",
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal valid file check", func() {
				err := json.Unmarshal(checksdConfigJson, checksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(checksdConfig.Checks).To(HaveLen(1))
				Expect(checksdConfig.Checks[0]).To(Equal(&checks.FileCheck{
					Name:      "name",
					Path:      "path",
					IfChanged: "ifchanged",
					Group:     "group",
					DependsOn: "dependson",
				}))

			})
		})

		Context("with process check config", func() {
			BeforeEach(func() {

				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "process",
						CheckProperties: map[string]interface{}{
							"Name":         "name",
							"Pidfile":      "pidfile",
							"StartProgram": "startprogram",
							"StopProgram":  "stopprogram",
							"Group":        "group",
							"DependsOn":    "dependson",
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal valid process check", func() {
				err := json.Unmarshal(checksdConfigJson, checksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(checksdConfig.Checks).To(HaveLen(1))
				Expect(checksdConfig.Checks[0]).To(Equal(&checks.ProcessCheck{
					Name:         "name",
					Pidfile:      "pidfile",
					StartProgram: "startprogram",
					StopProgram:  "stopprogram",
					Group:        "group",
					DependsOn:    "dependson",
				}))

			})
		})

		Context("multiple different checks", func() {
			BeforeEach(func() {

				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "process",
						CheckProperties: map[string]interface{}{
							"Name":         "name",
							"Pidfile":      "pidfile",
							"StartProgram": "startprogram",
							"StopProgram":  "stopprogram",
							"Group":        "group",
							"DependsOn":    "dependson",
						},
					},
					{
						Type: "icmp",
						CheckProperties: map[string]interface{}{
							"Address": "icmp_address",
							"Timeout": 1 * time.Second,
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal every check", func() {
				err := json.Unmarshal(checksdConfigJson, checksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(checksdConfig.Checks).To(HaveLen(2))
				Expect(checksdConfig.Checks[0]).To(Equal(&checks.ProcessCheck{
					Name:         "name",
					Pidfile:      "pidfile",
					StartProgram: "startprogram",
					StopProgram:  "stopprogram",
					Group:        "group",
					DependsOn:    "dependson",
				}))
				Expect(checksdConfig.Checks[1]).To(Equal(&network.IcmpCheck{
					Address: "icmp_address",
					Timeout: 1 * time.Second,
				}))

			})

		})

		Context("multiple same checks", func() {
			BeforeEach(func() {

				checksdConfig.ChecksConfig = []CheckConfig{
					{
						Type: "icmp",
						CheckProperties: map[string]interface{}{
							"Address": "another_icmp_address",
							"Timeout": 2 * time.Second,
						},
					},
					{
						Type: "icmp",
						CheckProperties: map[string]interface{}{
							"Address": "icmp_address",
							"Timeout": 1 * time.Second,
						},
					},
				}
				var err error
				checksdConfigJson, err = json.Marshal(checksdConfig)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unmarshal every check", func() {
				err := json.Unmarshal(checksdConfigJson, checksdConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(checksdConfig.Checks).To(HaveLen(2))
				Expect(checksdConfig.Checks[0]).To(Equal(&network.IcmpCheck{
					Address: "another_icmp_address",
					Timeout: 2 * time.Second,
				}))
				Expect(checksdConfig.Checks[1]).To(Equal(&network.IcmpCheck{
					Address: "icmp_address",
					Timeout: 1 * time.Second,
				}))

			})

		})
	})
})