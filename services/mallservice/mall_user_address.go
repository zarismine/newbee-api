package mallservice

import (
	"errors"
	"newbee/global"
	"newbee/models/mall"
	"newbee/models/mall/request"
	"newbee/pkg/dates"

	"github.com/jinzhu/copier"
)

var MallUserAddressService = newMallUserAddressService()

func newMallUserAddressService() *mallUserAddressService {
	return &mallUserAddressService{}
}

type mallUserAddressService struct {
}

func (m *mallUserAddressService) GetAddressByUserId (userid int) (userAddress []mall.MallUserAddress) {
	if err := global.DB.Where("user_id = ? and is_deleted = 0", userid).Find(&userAddress).Error; err != nil {
		return []mall.MallUserAddress{}
	}
	return userAddress
}

func (m *mallUserAddressService) Take (where ...interface{}) *mall.MallUserAddress {
	ret := &mall.MallUserAddress{}
	if err := global.DB.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (m *mallUserAddressService) Delete (address *mall.MallUserAddress) error {
	return global.DB.Delete(address).Error
}

func (m *mallUserAddressService) GetAddressByAddressId (addressid int)  *mall.MallUserAddress {
	ret := &mall.MallUserAddress{}
	global.DB.Where("address_id = ?", addressid).First(ret)
	return ret
}

func (m *mallUserAddressService) GetAddressByToken(token string) ([]mall.MallUserAddress, error) { 
	usertoken,err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return []mall.MallUserAddress{},err
	}
	addresslist := m.GetAddressByUserId(usertoken.UserId)
	return addresslist, err
}

func (m *mallUserAddressService) AddUserAddress(token string,req *request.AddAddressParam) (err error) { 
	userToken,err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return err
	}
	var Address mall.MallUserAddress
	copier.Copy(&Address, req)
	Address.CreateTime = dates.NowTimestamp()
	Address.UpdateTime = dates.NowTimestamp()
	Address.UserId = userToken.UserId
	if req.DefaultFlag == 1 {
		defaultAddress := m.Take("user_id=? and default_flag =1 and is_deleted = 0", userToken.UserId)
		if defaultAddress != nil {
			defaultAddress.UpdateTime = dates.NowTimestamp()
			defaultAddress.DefaultFlag = 0
			err = global.DB.Save(&defaultAddress).Error
			if err != nil {
				return
			}
		}
	} 
	err = global.DB.Create(&Address).Error
	return err
}

func (m *mallUserAddressService) EditUserAddress(token string,req *request.UpdateAddressParam) (err error) { 
	userToken,err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return err
	}
	userAddress := new(mall.MallUserAddress)
	global.DB.Where("address_id = ?",req.AddressId).First(userAddress)
	if userAddress.UserId != userToken.UserId {
		return errors.New("无权限")
	}
	if req.DefaultFlag == 1 {
		defaultAddress := m.Take("user_id=? and default_flag =1 and is_deleted = 0", userToken.UserId)
		if defaultAddress != nil {
			defaultAddress.UpdateTime = dates.NowTimestamp()
			defaultAddress.DefaultFlag = 0
			err = global.DB.Save(&defaultAddress).Error
			if err != nil {
				return
			}
		}
	}
	copier.Copy(userAddress, req)
	userAddress.UpdateTime = dates.NowTimestamp()
	err = global.DB.Save(&userAddress).Error
	return err
}

func (m *mallUserAddressService) GetDefaultAddressByToken (token string) (*mall.MallUserAddress, error) {
	userToken,err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return nil, err
	}
	defaultAddress := m.Take("user_id = ? and default_flag = 1 and is_deleted = 0", userToken.UserId)
	if defaultAddress != nil {
		return defaultAddress, nil
	}
	return nil, nil
}

func (m *mallUserAddressService) DeleteByAddressId (token string,addressid int) error {
	userToken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return err
	}
	deleteAddress := m.Take("address_id = ? and is_deleted = 0", addressid)
	if userToken.UserId != deleteAddress.UserId {
		return errors.New("无权限删除")
	}
	if deleteAddress != nil {
		return global.DB.Model(deleteAddress).UpdateColumns(map[string]interface{}{
			"update_time" : dates.NowTimestamp(),
			"is_deleted"  : 1,
		}).Error
	}
	return errors.New("地址已删除或不存在")
}