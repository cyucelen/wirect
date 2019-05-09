// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/cyucelen/wirect/model"

// PacketDatabase is an autogenerated mock type for the PacketDatabase type
type PacketDatabase struct {
	mock.Mock
	CreatedPackets []model.Packet
}

// CreatePacket provides a mock function with given fields: packet
func (_m *PacketDatabase) CreatePacket(packet *model.Packet) error {
	ret := _m.Called(packet)

	_m.CreatedPackets = append(_m.CreatedPackets, *packet)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Packet) error); ok {
		r0 = rf(packet)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetPacketsBySniffer provides a mock function with given fields: snifferMAC
func (_m *PacketDatabase) GetPacketsBySniffer(snifferMAC string) []model.Packet {
	ret := _m.Called(snifferMAC)

	var r0 []model.Packet
	if rf, ok := ret.Get(0).(func(string) []model.Packet); ok {
		r0 = rf(snifferMAC)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Packet)
		}
	}

	return r0
}

// GetPacketsBySnifferBetweenDates provides a mock function with given fields: snifferMAC, from, until
func (_m *PacketDatabase) GetPacketsBySnifferBetweenDates(snifferMAC string, from int64, until int64) []model.Packet {
	ret := _m.Called(snifferMAC, from, until)

	var r0 []model.Packet
	if rf, ok := ret.Get(0).(func(string, int64, int64) []model.Packet); ok {
		r0 = rf(snifferMAC, from, until)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Packet)
		}
	}

	return r0
}

// GetPacketsBySnifferSince provides a mock function with given fields: snifferMAC, since
func (_m *PacketDatabase) GetPacketsBySnifferSince(snifferMAC string, since int64) []model.Packet {
	ret := _m.Called(snifferMAC, since)

	var r0 []model.Packet
	if rf, ok := ret.Get(0).(func(string, int64) []model.Packet); ok {
		r0 = rf(snifferMAC, since)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Packet)
		}
	}

	return r0
}

// GetUniqueMACCountBySnifferBetweenDates provides a mock function with given fields: snifferMAC, from, until
func (_m *PacketDatabase) GetUniqueMACCountBySnifferBetweenDates(snifferMAC string, from int64, until int64) int {
	ret := _m.Called(snifferMAC, from, until)

	var r0 int
	if rf, ok := ret.Get(0).(func(string, int64, int64) int); ok {
		r0 = rf(snifferMAC, from, until)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}
