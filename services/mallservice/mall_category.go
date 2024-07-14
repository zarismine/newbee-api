package mallservice

import (
	// "fmt"
	"newbee/global"
	"newbee/models/mall"
	"newbee/models/manage"
	// "reflect"
	// "github.com/jinzhu/copier"
)

var MallCategoryService = newMallCategoryService()

func newMallCategoryService() *mallCategoryService {
	return &mallCategoryService{}
}

type mallCategoryService struct {
	FL []mall.FirstLevelCategoryVOS
}

func (m *mallCategoryService) GetList() ([]mall.FirstLevelCategoryVOS, error) {
	if m.FL != nil {
		return m.FL, nil
	}
	Res, err := m.GetCategoriesForIndex()
	m.FL = Res
	return m.FL, err
}
func reflectfromCategory(FirstCategories []*manage.MallGoodsCategory,targetReq []mall.ThirdLevelCategoryVOS) {
	for i := 0;i < len(FirstCategories);i++ {
        targetReq[i] = mall.ThirdLevelCategoryVOS{
            CategoryId   :  FirstCategories[i].CategoryId,
            CategoryLevel:  FirstCategories[i].CategoryLevel,
            CategoryName :  FirstCategories[i].CategoryName,
			ParentId     :  FirstCategories[i].ParentId ,
        }
	}
}

func (m *mallCategoryService) GetCategoriesForIndex() ([]mall.FirstLevelCategoryVOS, error) {
	var Categories []*manage.MallGoodsCategory
	global.DB.Table("tb_newbee_mall_goods_category").Where("category_level = ? and is_deleted = 0", 1).Order("category_rank DESC").Find(&Categories)
	TempFirst := make([]mall.ThirdLevelCategoryVOS,len(Categories))
	reflectfromCategory(Categories, TempFirst)
	FL := make([]mall.FirstLevelCategoryVOS,len(Categories))
	for i := 0;i < len(TempFirst);i++ {
		global.DB.Table("tb_newbee_mall_goods_category").Where("category_level = ? and is_deleted = 0 and parent_id = ?", 2, TempFirst[i].CategoryId).Order("category_rank DESC").Find(&Categories)
		TempSecond := make([]mall.ThirdLevelCategoryVOS,len(Categories))
		reflectfromCategory(Categories, TempSecond)
		SL := make([]mall.SecondLevelCategoryVOS,len(Categories))
		for j := 0;j < len(TempSecond);j++ {
			global.DB.Table("tb_newbee_mall_goods_category").Where("category_level = ? and is_deleted = 0 and parent_id = ?", 3, TempSecond[j].CategoryId).Order("category_rank DESC").Find(&Categories)
			TempThird := make([]mall.ThirdLevelCategoryVOS,len(Categories))
			reflectfromCategory(Categories, TempThird)
			SL[j] = mall.SecondLevelCategoryVOS {
				ThirdLevelCategoryVOS : TempSecond[j],
				ThirdLevelCategory: TempThird,
			}
		}
		FL[i] = mall.FirstLevelCategoryVOS {
			ThirdLevelCategoryVOS: TempFirst[i],
			SecondLevelCategory: SL,
		}
	}
	return FL, nil
}