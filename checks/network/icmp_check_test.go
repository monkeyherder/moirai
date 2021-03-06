package network_test

import (
	. "github.com/monkeyherder/salus/checks/network"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IcmpCheck", func() {
	var icmpCheck IcmpCheck

	BeforeEach(func() {
		icmpCheck = IcmpCheck{
			Address: "google.com",
			Timeout: 5 * time.Second,
		}
	})
	Context("Given a valid address", func() {
		It("should not return an error", func() {
			_, err := icmpCheck.Run()
			Expect(err).ToNot(HaveOccurred())
		})

		Context("that has disabled ICMP", func() {
			BeforeEach(func() {
				icmpCheck.Address = "test.com"
			})
			It("should return an error", func() {
				_, err := icmpCheck.Run()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("read udp 0.0.0.0:0: i/o timeout"))
			})
		})
	})

	Context("Given a invalid address", func() {
		BeforeEach(func() {
			icmpCheck = IcmpCheck{
				Address: "foo.bar.foo.bar",
			}
		})

		It("should return an error", func() {
			_, err := icmpCheck.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("lookup foo.bar.foo.bar.* no such host"))
		})
	})

})
