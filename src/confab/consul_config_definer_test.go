package confab_test

import (
	"confab"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConsulConfigDefiner", func() {
	Describe("GenerateConfiguration", func() {
		var consulConfig confab.ConsulConfig

		BeforeEach(func() {
			consulConfig = confab.GenerateConfiguration(confab.Config{})
		})

		Describe("datacenter", func() {
			It("defaults to empty string", func() {
				Expect(consulConfig.Datacenter).To(Equal(""))
			})

			Context("when the `consul.agent.datacenter` property is set", func() {
				It("uses that value", func() {
					consulConfig = confab.GenerateConfiguration(confab.Config{
						Consul: confab.ConfigConsul{
							Agent: confab.ConfigConsulAgent{
								Datacenter: "my-datacenter",
							},
						},
					})
					Expect(consulConfig.Datacenter).To(Equal("my-datacenter"))
				})
			})
		})

		Describe("domain", func() {
			It("defaults to `cf.internal`", func() {
				Expect(consulConfig.Domain).To(Equal("cf.internal"))
			})
		})

		Describe("data_dir", func() {
			It("defaults to `/var/vcap/store/consul_agent`", func() {
				Expect(consulConfig.DataDir).To(Equal("/var/vcap/store/consul_agent"))
			})
		})

		Describe("log_level", func() {
			It("defaults to empty string", func() {
				Expect(consulConfig.LogLevel).To(Equal(""))
			})

			Context("when the `consul.agent.log_level` property is set", func() {
				It("uses that value", func() {
					consulConfig = confab.GenerateConfiguration(confab.Config{
						Consul: confab.ConfigConsul{
							Agent: confab.ConfigConsulAgent{
								LogLevel: "some-log-level",
							},
						},
					})
					Expect(consulConfig.LogLevel).To(Equal("some-log-level"))
				})
			})
		})

		Describe("node_name", func() {
			It("uses the job name and index as the value", func() {
				consulConfig = confab.GenerateConfiguration(confab.Config{
					Node: confab.ConfigNode{
						Name:  "node_name",
						Index: 0,
					},
				})
				Expect(consulConfig.NodeName).To(Equal("node-name-0"))
			})
		})

		Describe("server", func() {
			It("defaults to false", func() {
				Expect(consulConfig.Server).To(BeFalse())
			})

			Context("when the `consul.agent.mode property` is `server`", func() {
				It("sets the value to true", func() {
					consulConfig = confab.GenerateConfiguration(confab.Config{
						Consul: confab.ConfigConsul{
							Agent: confab.ConfigConsulAgent{
								Mode: "server",
							},
						},
					})
					Expect(consulConfig.Server).To(BeTrue())
				})
			})

			Context("when the `consul.agent.mode` property is not `server`", func() {
				It("sets the value to false", func() {
					consulConfig = confab.GenerateConfiguration(confab.Config{
						Consul: confab.ConfigConsul{
							Agent: confab.ConfigConsulAgent{
								Mode: "banana",
							},
						},
					})
					Expect(consulConfig.Server).To(BeFalse())
				})
			})
		})

		Describe("ports", func() {
			It("defaults to a struct containing port 53 for DNS", func() {
				Expect(consulConfig.Ports).To(Equal(confab.ConsulConfigPorts{
					DNS: 53,
				}))
			})
		})

		Describe("rejoin_after_leave", func() {
			It("defaults to true", func() {
				Expect(consulConfig.RejoinAfterLeave).To(BeTrue())
			})
		})

		Describe("retry_join", func() {
			It("defaults to an empty slice of strings", func() {
				Expect(consulConfig.RetryJoin).To(Equal([]string{}))
			})

			Context("when `consul.agent.servers.lan` has a list of servers", func() {
				It("uses those values", func() {
					consulConfig = confab.GenerateConfiguration(confab.Config{
						Consul: confab.ConfigConsul{
							Agent: confab.ConfigConsulAgent{
								Servers: confab.ConfigConsulAgentServers{
									LAN: []string{
										"first-server",
										"second-server",
										"third-server",
									},
								},
							},
						},
					})
					Expect(consulConfig.RetryJoin).To(Equal([]string{
						"first-server",
						"second-server",
						"third-server",
					}))
				})
			})
		})

		Describe("bind_addr", func() {
			It("defaults to an empty string", func() {
				Expect(consulConfig.BindAddr).To(Equal(""))
			})

			Context("when `node.external_ip` is provided", func() {
				It("uses those values", func() {
					consulConfig = confab.GenerateConfiguration(confab.Config{
						Node: confab.ConfigNode{
							ExternalIP: "0.0.0.0",
						},
					})
					Expect(consulConfig.BindAddr).To(Equal("0.0.0.0"))
				})
			})
		})

		Describe("disable_remote_exec", func() {
			It("defaults to true", func() {
				Expect(consulConfig.DisableRemoteExec).To(BeTrue())
			})
		})

		Describe("disable_update_check", func() {
			It("defaults to true", func() {
				Expect(consulConfig.DisableUpdateCheck).To(BeTrue())
			})
		})

		Describe("protocol", func() {
			It("defaults to 0", func() {
				Expect(consulConfig.Protocol).To(Equal(0))
			})

			Context("when `consul.agent.protocol_version` is specified", func() {
				It("uses that value", func() {
					consulConfig = confab.GenerateConfiguration(confab.Config{
						Consul: confab.ConfigConsul{
							Agent: confab.ConfigConsulAgent{
								ProtocolVersion: 21,
							},
						},
					})
					Expect(consulConfig.Protocol).To(Equal(21))
				})
			})
		})

		Describe("verify_outgoing", func() {
			Context("when `consul.require_ssl` is false", func() {
				It("is nil", func() {
					Expect(consulConfig.VerifyOutgoing).To(BeNil())
				})
			})

			Context("when `consul.require_ssl` is true", func() {
				It("is true", func() {
					consulConfig = confab.GenerateConfiguration(confab.Config{
						Consul: confab.ConfigConsul{
							RequireSSL: true,
						},
					})
					Expect(consulConfig.VerifyOutgoing).NotTo(BeNil())
					Expect(*consulConfig.VerifyOutgoing).To(BeTrue())
				})
			})
		})

		Describe("verify_incoming", func() {
			Context("when `consul.require_ssl` is false", func() {
				It("is nil", func() {
					Expect(consulConfig.VerifyIncoming).To(BeNil())
				})
			})

			Context("when `consul.require_ssl` is true", func() {
				It("is true", func() {
					consulConfig = confab.GenerateConfiguration(confab.Config{
						Consul: confab.ConfigConsul{
							RequireSSL: true,
						},
					})
					Expect(consulConfig.VerifyIncoming).NotTo(BeNil())
					Expect(*consulConfig.VerifyIncoming).To(BeTrue())
				})
			})
		})

		Describe("verify_server_hostname", func() {
			Context("when `consul.require_ssl` is false", func() {
				It("is nil", func() {
					Expect(consulConfig.VerifyServerHostname).To(BeNil())
				})
			})

			Context("when `consul.require_ssl` is true", func() {
				It("is true", func() {
					consulConfig = confab.GenerateConfiguration(confab.Config{
						Consul: confab.ConfigConsul{
							RequireSSL: true,
						},
					})
					Expect(consulConfig.VerifyServerHostname).NotTo(BeNil())
					Expect(*consulConfig.VerifyServerHostname).To(BeTrue())
				})
			})
		})

		Describe("ca_file", func() {
			Context("when `consul.require_ssl` is false", func() {
				It("is nil", func() {
					Expect(consulConfig.CAFile).To(BeNil())
				})
			})

			Context("when `consul.require_ssl` is true", func() {
				It("is the location of the ca file", func() {
					consulConfig = confab.GenerateConfiguration(confab.Config{
						Consul: confab.ConfigConsul{
							RequireSSL: true,
						},
					})
					Expect(consulConfig.CAFile).NotTo(BeNil())
					Expect(*consulConfig.CAFile).To(Equal("/var/vcap/jobs/consul_agent/config/certs/ca.crt"))
				})
			})
		})

		Describe("key_file", func() {
			Context("when `consul.require_ssl` is false", func() {
				It("is nil", func() {
					Expect(consulConfig.KeyFile).To(BeNil())
				})
			})

			Context("when `consul.require_ssl` is true", func() {
				Context("when `consul.agent.mode` is `server`", func() {
					It("is the location of the server.key file", func() {
						consulConfig = confab.GenerateConfiguration(confab.Config{
							Consul: confab.ConfigConsul{
								RequireSSL: true,
								Agent: confab.ConfigConsulAgent{
									Mode: "server",
								},
							},
						})
						Expect(consulConfig.KeyFile).NotTo(BeNil())
						Expect(*consulConfig.KeyFile).To(Equal("/var/vcap/jobs/consul_agent/config/certs/server.key"))
					})
				})

				Context("when `consul.agent.mode` is not `server`", func() {
					It("is the location of the agent.key file", func() {
						consulConfig = confab.GenerateConfiguration(confab.Config{
							Consul: confab.ConfigConsul{
								RequireSSL: true,
							},
						})
						Expect(consulConfig.KeyFile).NotTo(BeNil())
						Expect(*consulConfig.KeyFile).To(Equal("/var/vcap/jobs/consul_agent/config/certs/agent.key"))
					})
				})
			})
		})

		Describe("cert_file", func() {
			Context("when `consul.require_ssl` is false", func() {
				It("is nil", func() {
					Expect(consulConfig.CertFile).To(BeNil())
				})
			})

			Context("when `consul.require_ssl` is true", func() {
				Context("when `consul.agent.mode` is `server`", func() {
					It("is the location of the server.crt file", func() {
						consulConfig = confab.GenerateConfiguration(confab.Config{
							Consul: confab.ConfigConsul{
								RequireSSL: true,
								Agent: confab.ConfigConsulAgent{
									Mode: "server",
								},
							},
						})
						Expect(consulConfig.CertFile).NotTo(BeNil())
						Expect(*consulConfig.CertFile).To(Equal("/var/vcap/jobs/consul_agent/config/certs/server.crt"))
					})
				})

				Context("when `consul.agent.mode` is not `server`", func() {
					It("is the location of the agent.key file", func() {
						consulConfig = confab.GenerateConfiguration(confab.Config{
							Consul: confab.ConfigConsul{
								RequireSSL: true,
							},
						})
						Expect(consulConfig.CertFile).NotTo(BeNil())
						Expect(*consulConfig.CertFile).To(Equal("/var/vcap/jobs/consul_agent/config/certs/agent.crt"))
					})
				})
			})
		})

		Describe("encrypt", func() {
			Context("when `consul.require_ssl` is false", func() {
				It("is nil", func() {
					Expect(consulConfig.Encrypt).To(BeNil())
				})
			})

			Context("when `consul.require_ssl` is true", func() {
				Context("when `consul.encrypt_keys` is empty", func() {
					It("is nil", func() {
						consulConfig = confab.GenerateConfiguration(confab.Config{
							Consul: confab.ConfigConsul{
								RequireSSL: true,
							},
						})
						Expect(consulConfig.Encrypt).To(BeNil())
					})
				})

				Context("when `consul.encrypt_keys` is provided with a key", func() {
					It("is an encoded version of that key", func() {
						consulConfig = confab.GenerateConfiguration(confab.Config{
							Consul: confab.ConfigConsul{
								RequireSSL:  true,
								EncryptKeys: []string{"banana"},
							},
						})
						Expect(consulConfig.Encrypt).NotTo(BeNil())
						Expect(*consulConfig.Encrypt).To(Equal("enqzXBmgKOy13WIGsmUk+g=="))
					})
				})
			})
		})

		Describe("bootstrap_expect", func() {
			Context("when `consul.agent.mode` is not `server`", func() {
				It("is nil", func() {
					Expect(consulConfig.BootstrapExpect).To(BeNil())
				})
			})

			Context("when `consul.agent.mode` is `server`", func() {
				It("sets it to the number of servers in the cluster", func() {
					consulConfig = confab.GenerateConfiguration(confab.Config{
						Consul: confab.ConfigConsul{
							Agent: confab.ConfigConsulAgent{
								Mode: "server",
								Servers: confab.ConfigConsulAgentServers{
									LAN: []string{
										"first-server",
										"second-server",
										"third-server",
									},
								},
							},
						},
					})
					Expect(consulConfig.BootstrapExpect).NotTo(BeNil())
					Expect(*consulConfig.BootstrapExpect).To(Equal(3))
				})
			})
		})
	})
})
