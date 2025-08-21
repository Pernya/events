package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin/binding"

	//"github.com/appleboy/go-fcm"

	"github.com/gin-gonic/gin"
	"github.com/restream/reindexer/v5"
	_ "github.com/restream/reindexer/v5/bindings/cproto"

	//"google.golang.org/api/fcm/v1"

	"io"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Product_list []struct {
	Offer_id string  `json:"offer_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price,omitempty"`
}
type Finish_checkout struct {
	Client_id           string       `json:"client_id"`           //	@description	Идентификатор клиента.
	User_id             string       `json:"user_id"`             //	@description	Идентификатор пользователя.
	Session_id          string       `json:"session_id"`          //	@description	Идентификатор сессии.
	Contact_key         string       `json:"contact_key"`         //	@description	Ключ контакта.
	Datetime            string       `json:"datetime"`            //	@description	Дата и время завершения оформления заказа.
	Source              string       `json:"source"`              //	@description	Источник трафика.
	Medium              string       `json:"medium"`              //	@description	Канал трафика.
	Campaign            string       `json:"campaign"`            //	@description	Кампания.
	Content             string       `json:"content"`             //	@description	Контент.
	Keywords            string       `json:"keywords"`            //	@description	Ключевые слова.
	Term                string       `json:"term"`                //	@description	Условия поиска.
	Currency            string       `json:"currency"`            //	@description	Валюта транзакции.
	Value               float64      `json:"value"`               //	@description	Общая стоимость товаров.
	Quantity            int          `json:"quantity"`            //	@description	Общее количество товаров.
	Coupon_code         string       `json:"coupon_code"`         //	@description	Промокод, использованный при оформлении заказа.
	Product_list        Product_list `json:"product_list"`        //	@description	Список товаров.
	Order_type          string       `json:"order_type"`          //	@description	Тип заказа.
	Discount_value      float64      `json:"discount_value"`      //	@description	Сумма скидки.
	Sku_quantity        int          `json:"sku_quantity"`        //	@description	Общее количество SKU.
	Sku_value           float64      `json:"sku_value"`           //	@description	Общая стоимость SKU.
	Shipping_type       string       `json:"shipping_type"`       //	@description	Тип доставки.
	Shipping_cost       float64      `json:"shipping_cost"`       //	@description	Стоимость доставки.
	Single_package      bool         `json:"single_package"`      //	@description	Является ли заказ одной посылкой.
	Payment_type        string       `json:"payment_type"`        //	@description	Способ оплаты.
	Recipient_type      string       `json:"recipient_type"`      //	@description	Тип получателя.
	Coupon_value        float64      `json:"coupon_value"`        //	@description	Сумма скидки по промокоду.
	Bonuses_spent_value float64      `json:"bonuses_spent_value"` //	@description	Сумма использованных бонусов.
	Bonuses_added_value float64      `json:"bonuses_added_value"` //	@description	Сумма начисленных бонусов.
	Total_weight        float64      `json:"total_weight"`        //	@description	Общий вес товаров.
	Geo_region          string       `json:"geo_region"`          //	@description	Регион пользователя.
	Geo_city            string       `json:"geo_city"`            //	@description	Город пользователя.
	Page_path           string       `json:"page_path"`           //	@description	Путь страницы завершения оформления заказа.
	Page_location       string       `json:"page_location"`       //	@description	Местоположение страницы.
	Screen_name         string       `json:"screen_name"`         //	@description	Название экрана.
	Screen_classname    string       `json:"screen_classname"`    //	@description	Класс экрана.
	Is_expressdelivery  bool         `json:"is_expressdelivery"`  //	@description	Доступна ли экспресс-доставка.
	Page_title          string       `json:"page_title"`          //	@description	Заголовок страницы.
	Referrer            string       `json:"referrer"`            //	@description	Источник перехода на страницу.
	Event_source        string       `json:"event_source"`        //	@description	Источник события.
	Is_auth             bool         `json:"is_auth"`             //	@description	Авторизован ли пользователь.
	Cart_id             string       `json:"cart_id"`             //	@description	Идентификатор корзины.
	Device_id           string       `json:"device_id,omitempty"` //	@description	Идентификатор устройства.
	Is_avatar           bool         `json:"is_avatar,omitempty"` //	@description	Проверка на аватар
	Device              string       `json:"device,omitempty"`    //	@description	Устройство
	Device_type         string       `json:"device_type"`         //	@description	Тип устройства
	Os                  string       `json:"os"`                  //	@description	ОС устройства
	Browser             string       `json:"browser"`             //	@description	Браузер

}

type Add_payment_info struct {
	Client_id          string       `json:"client_id"`           //	@description	Идентификатор клиента.
	User_id            string       `json:"user_id"`             //	@description	Идентификатор пользователя.
	Session_id         string       `json:"session_id"`          //	@description	Идентификатор сессии.
	Contact_key        string       `json:"contact_key"`         //	@description	Ключ контакта.
	Datetime           string       `json:"datetime"`            //	@description	Дата и время добавления данных для оплаты.
	Source             string       `json:"source"`              //	@description	Источник трафика.
	Medium             string       `json:"medium"`              //	@description	Канал трафика.
	Campaign           string       `json:"campaign"`            //	@description	Кампания.
	Content            string       `json:"content"`             //	@description	Контент.
	Keywords           string       `json:"keywords"`            //	@description	Ключевые слова.
	Term               string       `json:"term"`                //	@description	Условия поиска.
	Currency           string       `json:"currency"`            //	@description	Валюта транзакции.
	Value              float64      `json:"value"`               //	@description	Общая стоимость товаров.
	Product_list       Product_list `json:"product_list"`        //	@description	Список продуктов.
	Quantity           int          `json:"quantity"`            //	@description	Общее количество товаров.
	Sku_quantity       int          `json:"sku_quantity"`        //	@description	Общее количество SKU.
	Sku_value          float64      `json:"sku_value"`           //	@description	Общая стоимость SKU.
	Order_type         string       `json:"order_type"`          //	@description	Тип заказа.
	Discount_value     float64      `json:"discount_value"`      //	@description	Сумма скидки.
	Geo_region         string       `json:"geo_region"`          //	@description	Регион пользователя.
	Geo_city           string       `json:"geo_city"`            //	@description	Город пользователя.
	Page_path          string       `json:"page_path"`           //	@description	Путь страницы добавления информации для оплаты.
	Page_location      string       `json:"page_location"`       //	@description	Местоположение страницы.
	Screen_name        string       `json:"screen_name"`         //	@description	Название экрана.
	Screen_classname   string       `json:"screen_classname"`    //	@description	Класс экрана.
	Is_expressdelivery bool         `json:"is_expressdelivery"`  //	@description	Доступна ли экспресс-доставка.
	Page_title         string       `json:"page_title"`          //	@description	Заголовок страницы.
	Referrer           string       `json:"referrer"`            //	@description	Источник перехода на страницу.
	Event_source       string       `json:"event_source"`        //	@description	Источник события.
	Is_auth            bool         `json:"is_auth"`             //	@description	Авторизован ли пользователь.
	Cart_id            string       `json:"cart_id"`             //	@description	Идентификатор корзины.
	Device_id          string       `json:"device_id,omitempty"` //	@description	Идентификатор устройства.
	Is_avatar          bool         `json:"is_avatar,omitempty"` //	@description	Проверка на аватар
	//Screen_name        string       `json:"screen_name,omitempty"`      //	@description	Название экрана приложения
	//Screen_classname   string       `json:"screen_classname,omitempty"` //	@description	Название класса МП
	Device      string `json:"device,omitempty"` //	@description	Устройство
	Device_type string `json:"device_type"`      //	@description	Тип устройства
	Os          string `json:"os"`               //	@description	ОС устройства
	Browser     string `json:"browser"`          //	@description	Браузер
}
type Double_opt_in struct {
	Client_id        string `json:"client_id"`                  //	@description	Идентификатор клиента.
	User_id          string `json:"user_id"`                    //	@description	Идентификатор пользователя.
	Session_id       string `json:"session_id"`                 //	@description	Идентификатор сессии.
	Contact_key      string `json:"contact_key"`                //	@description	Ключ контакта.
	Datetime         string `json:"datetime"`                   //	@description	Дата и время подтверждения.
	Email            string `json:"email"`                      //	@description	Электронная почта пользователя.
	Url              string `json:"url"`                        //	@description	URL страницы, на которой произошло подтверждение.
	Event_source     string `json:"event_source"`               //	@description	Источник события.
	Device_id        string `json:"device_id,omitempty"`        //	@description	Идентификатор устройства.
	Screen_name      string `json:"screen_name,omitempty"`      //	@description	Название экрана приложения
	Screen_classname string `json:"screen_classname,omitempty"` //	@description	Название класса МП
	Device           string `json:"device,omitempty"`           //	@description	Устройство
	Device_type      string `json:"device_type"`                //	@description	Тип устройства
	Os               string `json:"os"`                         //	@description	ОС устройства
	Browser          string `json:"browser"`                    //	@description	Браузер
}
type Opt_in_confirmed struct {
	Client_id        string `json:"client_id"`                  //	@description	Идентификатор клиента.
	User_id          string `json:"user_id"`                    //	@description	Идентификатор пользователя.
	Session_id       string `json:"session_id"`                 //	@description	Идентификатор сессии.
	Contact_key      string `json:"contact_key"`                //	@description	Ключ контакта.
	Datetime         string `json:"datetime"`                   //	@description	Дата и время подтверждения.
	Email            string `json:"email"`                      //	@description	Электронная почта пользователя.
	Url              string `json:"url"`                        //	@description	URL страницы, на которой произошло подтверждение.
	Event_source     string `json:"event_source"`               //	@description	Источник события.
	Device_id        string `json:"device_id,omitempty"`        //	@description	Идентификатор устройства.
	Screen_name      string `json:"screen_name,omitempty"`      //	@description	Название экрана приложения
	Screen_classname string `json:"screen_classname,omitempty"` //	@description	Название класса МП
	Device           string `json:"device,omitempty"`           //	@description	Устройство
	Device_type      string `json:"device_type"`                //	@description	Тип устройства
	Os               string `json:"os"`                         //	@description	ОС устройства
	Browser          string `json:"browser"`                    //	@description	Браузер
}
type Add_contact_info struct {
	Client_id          string       `json:"client_id"`           //	@description	Идентификатор клиента.
	User_id            string       `json:"user_id"`             //	@description	Идентификатор пользователя.
	Session_id         string       `json:"session_id"`          //	@description	Идентификатор сессии.
	Contact_key        string       `json:"contact_key"`         //	@description	Ключ контакта.
	Datetime           string       `json:"datetime"`            //	@description	Дата и время добавления контактной информации.
	Source             string       `json:"source"`              //	@description	Источник трафика.
	Medium             string       `json:"medium"`              //	@description	Канал трафика.
	Campaign           string       `json:"campaign"`            //	@description	Кампания.
	Content            string       `json:"content"`             //	@description	Контент.
	Keywords           string       `json:"keywords"`            //	@description	Ключевые слова.
	Term               string       `json:"term"`                //	@description	Условия поиска.
	Currency           string       `json:"currency"`            //	@description	Валюта транзакции.
	Value              float64      `json:"value"`               //	@description	Общая стоимость товаров.
	Product_list       Product_list `json:"product_list"`        //	@description	Список продуктов.
	Quantity           int          `json:"quantity"`            //	@description	Общее количество товаров.
	Sku_quantity       int          `json:"sku_quantity"`        //	@description	Общее количество SKU.
	Recipient_type     string       `json:"recipient_type"`      //	@description	Тип получателя
	Sku_value          float64      `json:"sku_value"`           //	@description	Общая стоимость SKU.
	Order_type         string       `json:"order_type"`          //	@description	Тип заказа.
	Discount_value     float64      `json:"discount_value"`      //	@description	Сумма скидки.
	Geo_region         string       `json:"geo_region"`          //	@description	Регион пользователя.
	Geo_city           string       `json:"geo_city"`            //	@description	Город пользователя.
	Page_path          string       `json:"page_path"`           //	@description	Путь страницы добавления контактной информации.
	Page_location      string       `json:"page_location"`       //	@description	Местоположение страницы.
	Screen_name        string       `json:"screen_name"`         //	@description	Название экрана.
	Screen_classname   string       `json:"screen_classname"`    //	@description	Класс экрана.
	Is_expressdelivery bool         `json:"is_expressdelivery"`  //	@description	Доступна ли экспресс-доставка.
	Page_title         string       `json:"page_title"`          //	@description	Заголовок страницы.
	Referrer           string       `json:"referrer"`            //	@description	Источник перехода на страницу.
	Event_source       string       `json:"event_source"`        //	@description	Источник события.
	Is_auth            bool         `json:"is_auth"`             //	@description	Авторизован ли пользователь.
	Cart_id            string       `json:"cart_id"`             //	@description	Идентификатор корзины.
	Device_id          string       `json:"device_id,omitempty"` //	@description	Идентификатор устройства.
	Is_avatar          bool         `json:"is_avatar,omitempty"` //	@description	Проверка на аватар
	Device             string       `json:"device,omitempty"`    //	@description	Устройство
	Device_type        string       `json:"device_type"`         //	@description	Тип устройства
	Os                 string       `json:"os"`                  //	@description	ОС устройства
	Browser            string       `json:"browser"`             //	@description	Браузер

}
type View_Cart struct {
	Client_id               string       `json:"client_id,omitempty"`        //	@description	Идентификатор клиента.
	User_id                 string       `json:"user_id,omitempty"`          //	@description	Идентификатор пользователя.
	Session_id              string       `json:"session_id,omitempty"`       //	@description	Идентификатор сессии.
	Contact_key             string       `json:"contact_key,omitempty"`      //	@description	Идентификатор из Манзана
	Source                  string       `json:"source"`                     //	@description	Источник трафика.
	Medium                  string       `json:"medium"`                     //	@description	Канал трафика.
	Campaign                string       `json:"сampaign"`                   //	@description	Кампания.
	Content                 string       `json:"content"`                    //	@description	Контент.
	Datetime                string       `json:"datetime"`                   //	@decription Дата и время события
	Keywords                string       `json:"keywords"`                   //	@description	Ключевые слова.
	Term                    string       `json:"term"`                       //	@description	Условия поиска.
	Value                   float64      `json:"value"`                      //	@description	Общая стоимость товаров в корзине.
	Currency                string       `json:"currency"`                   //	@description	Валюта транзакции.
	Product_list            Product_list `json:"product_list"`               //	@description	Список продуктов в корзине.
	Quantity                int          `json:"quantity"`                   //	@description	Общее количество товаров.
	Sku_quantity            float64      `json:"sku_quantity"`               //	@description	Общее количество SKU в корзине.
	Sku_value               float64      `json:"sku_value"`                  //	@description	Общая стоимость SKU в корзине.
	Discount_value          float64      `json:"discount_value"`             //	@description	Сумма скидки на товары в корзине.
	Cart_type               string       `json:"cart_type"`                  //	@description	Тип корзины (например, временная, постоянная).
	Geo_region              string       `json:"geo_region"`                 //	@description	Регион пользователя.
	Geo_city                string       `json:"geo_city"`                   //	@description	Город пользователя.
	Page_path               string       `json:"page_path"`                  //	@description	Путь страницы просмотра корзины.
	Page_location           string       `json:"page_location"`              //	@description	Местоположение страницы.
	Is_expressdelivery      bool         `json:"is_expressdelivery"`         //	@description	Доступна ли экспресс-доставка.
	Page_title              string       `json:"page_title"`                 //	@description	Заголовок страницы.
	Referrer                string       `json:"referrer"`                   //	@description	Источник перехода на страницу.
	Is_subscriptiondelivery bool         `json:"is_subscriptiondelivery"`    //	@description	Доступна ли подписка на доставку.
	Event_source            string       `json:"event_source"`               //	@description	Источник события.
	Is_auth                 bool         `json:"is_auth"`                    //	@description	Авторизован ли пользователь.
	Cart_id                 string       `json:"cart_id"`                    //	@description	Идентификатор корзины.
	Device_id               string       `json:"device_id,omitempty"`        //	@description	Идентификатор устройства.
	Screen_name             string       `json:"screen_name,omitempty"`      //	@description	Название экрана приложения
	Screen_classname        string       `json:"screen_classname,omitempty"` //	@description	Название класса МП
	Is_avatar               bool         `json:"is_avatar,omitempty"`        //	@description	Проверка на аватар
	Device                  string       `json:"device,omitempty"`           //	@description	Устройство
	Device_type             string       `json:"device_type"`                //	@description	Тип устройства
	Os                      string       `json:"os"`                         //	@description	ОС устройства
	Browser                 string       `json:"browser"`                    //	@description	Браузер
}
type Add_to_Wishlist struct {
	Client_id               string  `json:"client_id,omitempty"`        //	@description	Идентификатор клиента.
	User_id                 string  `json:"user_id,omitempty"`          //	@description	Идентификатор пользователя.
	Session_id              string  `json:"session_id,omitempty"`       //	@description	Идентификатор сессии.
	Contact_key             string  `json:"contact_key,omitempty"`      //	@description	Идентификатор из Манзана
	Source                  string  `json:"source"`                     //	@description	Источник трафика.
	Medium                  string  `json:"medium"`                     //	@description	Канал трафика.
	Campaign                string  `json:"сampaign"`                   //	@description	Кампания.
	Content                 string  `json:"content"`                    //	@description	Контент.
	Keywords                string  `json:"keywords"`                   //	@description	Ключевые слова.
	Datetime                string  `json:"datetime"`                   //	@decription Дата и время события
	Term                    string  `json:"term"`                       //	@description	Условия поиска.
	Product_id              string  `json:"product_id"`                 //	@description	Идентификатор товара.
	Value                   float64 `json:"value"`                      //	@description	Стоимость товара.
	Search_term             string  `json:"search_term"`                //	@description	Поисковый запрос.
	Offer_id                string  `json:"offer_id"`                   //	@description	Идентификатор предложения.
	Name                    string  `json:"name"`                       //	@description	Название товара.
	Brand                   string  `json:"brand"`                      //	@description	Бренд товара.
	Currency                string  `json:"currency"`                   //	@description	Валюта транзакции.
	Discount_value          float64 `json:"discount_value"`             //	@description	Сумма скидки на товар.
	Index                   float64 `json:"index"`                      //	@description	Индекс продукта в списке.
	Item_list_id            string  `json:"item_list_id"`               //	@description	Идентификатор списка товаров.
	Item_list_name          string  `json:"item_list_name"`             //	@description	Название списка товаров.
	Category_id             int     `json:"category_id"`                //	@description	Идентификатор категории.
	Category_name           string  `json:"category_name"`              //	@description	Название категории.
	Category_hierarchy      string  `json:"category_hierarchy"`         //	@description	Иерархия категорий.
	Badge                   string  `json:"badge"`                      //	@description	Значок товара.
	Promo                   string  `json:"promo"`                      //	@description	Промо для товара.
	Rating                  float64 `json:"rating"`                     //	@description	Рейтинг товара.
	Availability            bool    `json:"availability"`               //	@description	Наличие товара.
	Geo_region              string  `json:"geo_region"`                 //	@description	Регион пользователя.
	Geo_city                string  `json:"geo_city"`                   //	@description	Город пользователя.
	Page_path               string  `json:"page_path"`                  //	@description	Путь страницы.
	Page_location           string  `json:"page_location"`              //	@description	Местоположение страницы.
	Is_expressdelivery      bool    `json:"is_expressdelivery"`         //	@description	Доступна ли экспресс-доставка.
	Page_title              string  `json:"page_title"`                 //	@description	Заголовок страницы.
	Referrer                string  `json:"referrer"`                   //	@description	Источник перехода на страницу.
	Is_subscriptiondelivery bool    `json:"is_subscriptiondelivery"`    //	@description	Доступна ли подписка на доставку.
	Event_source            string  `json:"event_source"`               //	@description	Источник события.
	Is_auth                 bool    `json:"is_auth"`                    //	@description	Авторизован ли пользователь.
	Cart_id                 string  `json:"cart_id"`                    //	@description	Идентификатор корзины.
	Device_id               string  `json:"device_id,omitempty"`        //	@description	Идентификатор устройства.
	Screen_name             string  `json:"screen_name,omitempty"`      //	@description	Название экрана приложения
	Screen_classname        string  `json:"screen_classname,omitempty"` //	@description	Название класса МП
	Device                  string  `json:"device,omitempty"`           //	@description	Устройство
	Device_type             string  `json:"device_type"`                //	@description	Тип устройства
	Os                      string  `json:"os"`                         //	@description	ОС устройства
	Browser                 string  `json:"browser"`                    //	@description	Браузер
}

type Remove_from_cart struct {
	Client_id               string       `json:"client_id,omitempty"`     //	@description	Идентификатор клиента.
	User_id                 string       `json:"user_id,omitempty"`       //	@description	Идентификатор пользователя.
	Session_id              string       `json:"session_id,omitempty"`    //	@description	Идентификатор сессии.
	Contact_key             string       `json:"contact_key,omitempty"`   //	@description	Идентификатор из Манзана
	Source                  string       `json:"source"`                  //	@description	Источник трафика.
	Medium                  string       `json:"medium"`                  //	@description	Канал трафика.
	Campaign                string       `json:"сampaign"`                //	@description	Кампания.
	Content                 string       `json:"content"`                 //	@description	Контент.
	Datetime                string       `json:"datetime"`                //	@decription Дата и время события
	Keywords                string       `json:"keywords"`                //	@description	Ключевые слова.
	Term                    string       `json:"term"`                    //	@description	Условия поиска.
	Value                   float64      `json:"value"`                   //	@description	Стоимость удалённого товара.
	Currency                string       `json:"currency"`                //	@description	Валюта транзакции.
	Product_id              string       `json:"product_id"`              //	@description	Идентификатор товара.
	Quantity                int          `json:"quantity"`                //	@description	Количество удалённых товаров.
	Cart_id                 string       `json:"cart_id"`                 //	@description	Идентификатор корзины.
	Search_term             string       `json:"search_term"`             //	@description	Поисковый запрос.
	Offer_id                string       `json:"offer_id"`                //	@description	Идентификатор предложения.
	Discount_value          float64      `json:"discount_value"`          //	@description	Сумма скидки на товар.
	Name                    string       `json:"name"`                    //	@description	Название товара.
	Brand                   string       `json:"brand"`                   //	@description	Бренд товара.
	Badge                   string       `json:"badge"`                   //	@description	Значок товара.
	Promo                   string       `json:"promo"`                   //	@description	Промо для товара.
	Rating                  float64      `json:"rating"`                  //	@description	Рейтинг товара.
	Item_list_id            string       `json:"item_list_id"`            //	@description	Идентификатор списка товаров.
	Item_list_name          string       `json:"item_list_name"`          //	@description	Название списка товаров.
	Category_id             int          `json:"category_id"`             //	@description	Идентификатор категории.
	Category_name           string       `json:"category_name"`           //	@description	Название категории.
	Category_hierarchy      string       `json:"category_hierarchy"`      //	@description	Иерархия категорий.
	Product_list            Product_list `json:"product_list"`            //	@description	Идентификаторы товаров в категории.
	Active_filters          interface{}  `json:"active_filters"`          //	@description	Активные фильтры, применённые пользователем.
	Geo_region              string       `json:"geo_region"`              //	@description	Регион пользователя.
	Geo_city                string       `json:"geo_city"`                //	@description	Город пользователя.
	Page_path               string       `json:"page_path"`               //	@description	Путь страницы.
	Page_location           string       `json:"page_location"`           //	@description	Местоположение страницы.
	Is_expressdelivery      bool         `json:"is_expressdelivery"`      //	@description	Доступна ли экспресс-доставка.
	Page_title              string       `json:"page_title"`              //	@description	Заголовок страницы.
	Referrer                string       `json:"referrer"`                //	@description	Источник перехода на страницу.
	Is_subscriptiondelivery bool         `json:"is_subscriptiondelivery"` //	@description	Доступна ли подписка на доставку.
	Event_source            string       `json:"event_source"`            //	@description	Источник события.
	Is_auth                 bool         `json:"is_auth"`                 //	@description	Авторизован ли пользователь.
	Device_id               string       `json:"device_id,omitempty"`     //	@description	Идентификатор устройства.
	Screen_name             string       `json:"screen_name,omitempty"`   //	@description	Название экрана приложения
	Screen_classname        string       `json:"screen_classname"`        //	@description	Название класса МП
	Is_avatar               bool         `json:"is_avatar,omitempty"`     //	@description	Проверка на аватар
	Device                  string       `json:"device,omitempty"`        //	@description	Устройство
	Device_type             string       `json:"device_type"`             //	@description	Тип устройства
	Os                      string       `json:"os"`                      //	@description	ОС устройства
	Browser                 string       `json:"browser"`                 //	@description	Браузер
}
type Add_shipping_info struct {
	ClientID          string       `json:"client_id"`           // @description Unique identifier for a user associated with device and browser.
	Device_id         string       `json:"device_id,omitempty"` //	@description	Идентификатор устройства.
	UserID            string       `json:"user_id"`             // @description Unique identifier for a registered user.
	SessionID         string       `json:"session_id"`          // @description Identifier for the user's session.
	ContactKey        string       `json:"contact_key"`         // @description Contact identifier from an external system.
	Datetime          string       `json:"datetime"`            // @description Date and time of the event.
	Source            string       `json:"source"`              // @description Source parameter linked to traffic source.
	Medium            string       `json:"medium"`              // @description Medium parameter linked to traffic source.
	Campaign          string       `json:"campaign"`            // @description Campaign parameter linked to traffic source.
	Content           string       `json:"content"`             // @description Content parameter linked to traffic source.
	Keywords          string       `json:"keywords"`            // @description Keywords parameter linked to traffic source.
	Term              string       `json:"term"`                // @description Term parameter linked to traffic source.
	Currency          string       `json:"currency"`            // @description Currency used for the transaction.
	Value             float64      `json:"value"`               // @description Value of the transaction.
	ProductList       Product_list `json:"product_list"`        // @description List of products involved in the transaction.
	ShippingType      string       `json:"shipping_type"`       // @description Type of shipping method used.
	ShippingCost      float64      `json:"shipping_cost"`       // @description Cost of shipping.
	SinglePackage     bool         `json:"single_package"`      // @description Indicates if the order is shipped as a single package.
	OrderType         string       `json:"order_type"`          // @description Type of order placed.
	GeoRegion         string       `json:"geo_region"`          // @description Geographic region for delivery.
	GeoCity           string       `json:"geo_city"`            // @description City for delivery.
	PagePath          string       `json:"page_path"`           // @description Path of the page where the transaction was initiated.
	PageLocation      string       `json:"page_location"`       // @description URL location of the page where the transaction was initiated.
	ScreenName        string       `json:"screen_name"`         // @description Name of the screen where the transaction was initiated.
	ScreenClassname   string       `json:"screen_classname"`    // @description Class name of the screen where the transaction was initiated.
	IsExpressDelivery bool         `json:"is_expressdelivery"`  // @description Indicates if the delivery is express.
	PageTitle         string       `json:"page_title"`          // @description Title of the page where the transaction was initiated.
	Referrer          string       `json:"referrer"`            // @description Referrer URL.
	EventSource       string       `json:"event_source"`        // @description Source of the event.
	IsAuth            bool         `json:"is_auth"`             // @description Indicates if the user is authenticated.
	CartID            string       `json:"cart_id"`             // @description Identifier for the user's cart.
	Is_avatar         bool         `json:"is_avatar,omitempty"` //	@description	Проверка на аватар
	Device            string       `json:"device,omitempty"`    //	@description	Устройство
	Device_type       string       `json:"device_type"`         //	@description	Тип устройства
	Os                string       `json:"os"`                  //	@description	ОС устройства
	Browser           string       `json:"browser"`             //	@description	Браузер

}

// Purchase_item содержит информацию о доставке и заказе для каждого купленного товара.
// @Description Содержит информацию о доставке и заказе для каждого купленного товара.
// @Tags         purchase_item
type Purchase_item struct {
	ClientID            string  `json:"client_id"`                  // @description Уникальный идентификатор пользователя.
	Device_id           string  `json:"device_id,omitempty"`        //	@description	Идентификатор устройства.
	UserID              string  `json:"user_id"`                    // @description Уникальный идентификатор зарегистрированного пользователя.
	SessionID           string  `json:"session_id"`                 // @description Идентификатор сессии пользователя.
	ContactKey          string  `json:"contact_key"`                // @description Идентификатор контакта.
	Datetime            string  `json:"datetime"`                   // @description Дата и время события.
	Source              string  `json:"source"`                     // @description Источник трафика.
	Medium              string  `json:"medium"`                     // @description Канал трафика.
	Campaign            string  `json:"campaign"`                   // @description Название рекламной кампании.
	Content             string  `json:"content"`                    // @description Содержание кампании.
	Keywords            string  `json:"keywords"`                   // @description Ключевые слова.
	Term                string  `json:"term"`                       // @description Условия поиска.
	GeoCity             string  `json:"geo_city"`                   // @description Город доставки.
	GeoRegion           string  `json:"geo_region"`                 // @description Регион доставки.
	CartID              string  `json:"cart_id"`                    // @description Идентификатор корзины пользователя.
	PackageID           string  `json:"package_id"`                 // @description Идентификатор упаковки.
	UniqueTransactionID string  `json:"unique_transaction_id"`      // @description Уникальный идентификатор транзакции.
	OrderType           string  `json:"order_type"`                 // @description Тип заказа.
	AvatarOrder         bool    `json:"avatar_order,omitempty"`     // @description Заказ оформлен аватаром.
	PurchaseSourceType  string  `json:"purchase_source_type"`       // @description Источник покупки.
	ProductID           string  `json:"product_id"`                 // @description Идентификатор товара.
	Brand               string  `json:"brand"`                      // @description Бренд товара.
	Title               string  `json:"title"`                      // @description Название товара.
	OfferID             string  `json:"offer_id"`                   // @description Идентификатор предложения.
	ProductSTM          bool    `json:"product_stm"`                // @description Товар собственной торговой марки.
	ItemPrice           float64 `json:"item_price"`                 // @description Цена товара.
	OriginalPrice       float64 `json:"original_price"`             // @description Оригинальная цена товара.
	DiscountPercentage  float64 `json:"discount_percentage"`        // @description Процент скидки.
	DiscountValue       float64 `json:"discount_value"`             // @description Сумма скидки.
	Quantity            int     `json:"quantity"`                   // @description Количество товаров.
	TotalPrice          float64 `json:"total_price"`                // @description Общая стоимость товаров после скидки.
	ProductFood         bool    `json:"product_food"`               // @description Продукт питания.
	ProductSpec         string  `json:"product_spec"`               // @description Спецификация продукта.
	ProductWeight       float64 `json:"product_weight"`             // @description Вес товара.
	ProductWearSize     string  `json:"product_wear_size"`          // @description Размер изделия.
	ProductColor        string  `json:"product_color"`              // @description Цвет товара.
	ProductTaste        string  `json:"product_taste"`              // @description Вкус товара.
	ProductFarma        bool    `json:"product_farma"`              // @description Ветеринарный продукт.
	ProductFarmaType    string  `json:"product_farma_type"`         // @description Тип ветеринарного продукта.
	PetAge              string  `json:"pet_age"`                    // @description Возраст питомца.
	PetSize             string  `json:"pet_size"`                   // @description Размер питомца.
	PetType             string  `json:"pet_type"`                   // @description Тип питомца.
	EventSource         string  `json:"event_source"`               // @description Источник события.
	Category1           string  `json:"category_1"`                 // @description Название категории первого уровня.
	Category1ID         string  `json:"category_1_id"`              // @description Идентификатор категории первого уровня.
	Category2           string  `json:"category_2"`                 // @description Название категории второго уровня.
	Category2ID         string  `json:"category_2_id"`              // @description Идентификатор категории второго уровня.
	Category3           string  `json:"category_3"`                 // @description Название категории третьего уровня.
	Category3ID         string  `json:"category_3_id"`              // @description Идентификатор категории третьего уровня.
	Category4           string  `json:"category_4"`                 // @description Название категории четвертого уровня.
	Category4ID         string  `json:"category_4_id"`              // @description Идентификатор категории четвертого уровня.
	Category5           string  `json:"category_5"`                 // @description Название категории пятого уровня.
	Category5ID         string  `json:"category_5_id"`              // @description Идентификатор категории пятого уровня.
	Screen_name         string  `json:"screen_name,omitempty"`      //	@description	Название экрана приложения
	Screen_classname    string  `json:"screen_classname,omitempty"` //	@description	Название класса МП
	Is_avatar           bool    `json:"is_avatar,omitempty"`        //	@description	Проверка на аватар
	Device              string  `json:"device,omitempty"`           //	@description	Устройство
	Device_type         string  `json:"device_type"`                //	@description	Тип устройства
	Os                  string  `json:"os"`                         //	@description	ОС устройства
	Browser             string  `json:"browser"`                    //	@description	Браузер
}

// Express_delivery содержит информацию о доставке заказов в экспресс-режиме.
// @Description Содержит информацию о доставке заказов в экспресс-режиме.
// @Tags         express_delivery
type Express_delivery struct {
	ClientID          string `json:"client_id"`           // @description Уникальный идентификатор пользователя.
	DeviceID          string `json:"device_id"`           // @description Уникальный идентификатор устройства пользователя.
	UserID            string `json:"user_id"`             // @description Уникальный идентификатор зарегистрированного пользователя.
	SessionID         string `json:"session_id"`          // @description Идентификатор сессии пользователя.
	ContactKey        string `json:"contact_key"`         // @description Идентификатор контакта.
	Datetime          string `json:"datetime"`            // @description Дата и время события.
	Source            string `json:"source"`              // @description Источник трафика.
	Medium            string `json:"medium"`              // @description Канал трафика.
	Campaign          string `json:"campaign"`            // @description Название рекламной кампании.
	Content           string `json:"content"`             // @description Содержание кампании.
	Keywords          string `json:"keywords"`            // @description Ключевые слова.
	Term              string `json:"term"`                // @description Условия поиска.
	Action            string `json:"action"`              // @description Действие пользователя.
	GeoRegion         string `json:"geo_region"`          // @description Регион доставки.
	GeoCity           string `json:"geo_city"`            // @description Город доставки.
	PagePath          string `json:"page_path"`           // @description Путь к странице.
	PageLocation      string `json:"page_location"`       // @description URL-адрес страницы.
	ScreenName        string `json:"screen_name"`         // @description Название экрана.
	ScreenClassname   string `json:"screen_classname"`    // @description Класс экрана.
	IsExpressDelivery bool   `json:"is_expressdelivery"`  // @description Признак экспресс-доставки.
	PageTitle         string `json:"page_title"`          // @description Заголовок страницы.
	Referrer          string `json:"referrer"`            // @description Реферер.
	EventSource       string `json:"event_source"`        // @description Источник события.
	IsAuth            bool   `json:"is_auth"`             // @description Аутентификация пользователя.
	CartID            string `json:"cart_id"`             // @description Идентификатор корзины пользователя.
	Is_avatar         bool   `json:"is_avatar,omitempty"` //	@description	Проверка на аватар
	Device            string `json:"device,omitempty"`    //	@description	Устройство
	Device_type       string `json:"device_type"`         //	@description	Тип устройства
	Os                string `json:"os"`                  //	@description	ОС устройства
	Browser           string `json:"browser"`             //	@description	Браузер

}

// Comment содержит информацию о комментариях и оценках пользователей для товаров.
// @Description Содержит информацию о комментариях и оценках пользователей для товаров.
// @Tags         comment
type Comment struct {
	ClientID         string  `json:"client_id"`                  // @description Уникальный идентификатор пользователя.
	DeviceID         string  `json:"device_id"`                  // @description Уникальный идентификатор устройства пользователя.
	UserID           string  `json:"user_id"`                    // @description Уникальный идентификатор зарегистрированного пользователя.
	SessionID        string  `json:"session_id"`                 // @description Идентификатор сессии пользователя.
	ContactKey       string  `json:"contact_key"`                // @description Идентификатор контакта.
	Datetime         string  `json:"datetime"`                   // @description Дата и время события.
	Source           string  `json:"source"`                     // @description Источник трафика.
	Medium           string  `json:"medium"`                     // @description Канал трафика.
	Campaign         string  `json:"campaign"`                   // @description Название рекламной кампании.
	Content          string  `json:"content"`                    // @description Содержание кампании.
	Keywords         string  `json:"keywords"`                   // @description Ключевые слова.
	Term             string  `json:"term"`                       // @description Условия поиска.
	ProductID        string  `json:"product_id"`                 // @description Идентификатор товара.
	EventSource      string  `json:"event_source"`               // @description Источник события.
	OfferID          string  `json:"offer_id"`                   // @description Идентификатор предложения.
	Title            string  `json:"title"`                      // @description Название товара.
	Brand            string  `json:"brand"`                      // @description Бренд товара.
	CommentID        string  `json:"comment_id"`                 // @description Идентификатор комментария.
	FotoCount        int     `json:"foto_count"`                 // @description Количество фотографий в отзыве.
	Rating           float64 `json:"rating"`                     // @description Оценка товара.
	ReviewText       string  `json:"review_text"`                // @description Текст отзыва.
	Category1        string  `json:"category_1"`                 // @description Название категории первого уровня.
	Category1ID      string  `json:"category_1_id"`              // @description Идентификатор категории первого уровня.
	Category2        string  `json:"category_2"`                 // @description Название категории второго уровня.
	Category2ID      string  `json:"category_2_id"`              // @description Идентификатор категории второго уровня.
	Category3        string  `json:"category_3"`                 // @description Название категории третьего уровня.
	Category3ID      string  `json:"category_3_id"`              // @description Идентификатор категории третьего уровня.
	Category4        string  `json:"category_4"`                 // @description Название категории четвертого уровня.
	Category4ID      string  `json:"category_4_id"`              // @description Идентификатор категории четвертого уровня.
	Category5        string  `json:"category_5"`                 // @description Название категории пятого уровня.
	Category5ID      string  `json:"category_5_id"`              // @description Идентификатор категории пятого уровня.
	Screen_name      string  `json:"screen_name,omitempty"`      //	@description	Название экрана приложения
	Screen_classname string  `json:"screen_classname,omitempty"` //	@description	Название класса МП
	Is_avatar        bool    `json:"is_avatar,omitempty"`        //	@description	Проверка на аватар
	Device           string  `json:"device,omitempty"`           //	@description	Устройство
	Device_type      string  `json:"device_type"`                //	@description	Тип устройства
	Os               string  `json:"os"`                         //	@description	ОС устройства
	Browser          string  `json:"browser"`                    //	@description	Браузер
}

// Search содержит информацию о поисковых запросах пользователей.
// @Description Содержит информацию о поисковых запросах пользователей.
// @Tags         search
type Search struct {
	Client_id          string `json:"client_id,omitempty"`
	Device_id          string `json:"device_id,omitempty"`
	Session_id         string `json:"session_id"`
	Contact_key        string `json:"contact_key"`
	User_id            string `json:"user_id"`
	DateTime           string `json:"date_time"`
	Source             string `json:"source"`
	Medium             string `json:"medium"`
	Campaign           string `json:"campaign"`
	Content            string `json:"content"`
	Keywords           string `json:"keywords"`
	Term               string `json:"term"`
	Search_term        string `json:"search_term"`
	Action             string `json:"action"`
	Suggest_term       string `json:"suggest_term"`
	Geo_region         string `json:"geo_region"`
	Geo_city           string `json:"geo_city"`
	Page_path          string `json:"page_path"`
	Page_location      string `json:"page_location"`
	Screen_name        string `json:"screen_name"`
	Screen_classname   string `json:"screen_classname"`
	Is_expressdelivery bool   `json:"is_expressdelivery"`
	Page_title         string `json:"page_title"`
	Referrer           string `json:"referrer"`
	Event_source       string `json:"event_source"`
	Is_auth            bool   `json:"is_auth"`
	Cart_id            string `json:"cart_id"`
	Is_avatar          bool   `json:"is_avatar,omitempty"` //	@description	Проверка на аватар
	Device             string `json:"device,omitempty"`    //	@description	Устройство
	Device_type        string `json:"device_type"`         //	@description	Тип устройства
	Os                 string `json:"os"`                  //	@description	ОС устройства
	Browser            string `json:"browser"`             //	@description	Браузер
}

// Notification содержит информацию о уведомлениях, отправляемых пользователям.
// @Description Содержит информацию о уведомлениях, отправляемых пользователям.
// @Tags         notification
//
//	type Notification struct {
//		ClientID         string       `json:"client_id"`                  // @description Уникальный идентификатор пользователя.
//		Device_id        string       `json:"device_id,omitempty"`        //	@description	Идентификатор устройства.
//		UserID           string       `json:"user_id"`                    // @description Уникальный идентификатор зарегистрированного пользователя.
//		SessionID        string       `json:"session_id"`                 // @description Идентификатор сессии пользователя.
//		ContactKey       string       `json:"contact_key"`                // @description Идентификатор контакта.
//		Datetime         string       `json:"datetime"`                   // @description Дата и время события.
//		Source           string       `json:"source"`                     // @description Источник трафика.
//		Medium           string       `json:"medium"`                     // @description Канал трафика.
//		Campaign         string       `json:"campaign"`                   // @description Название рекламной кампании.
//		Content          string       `json:"content"`                    // @description Содержание кампании.
//		Keywords         string       `json:"keywords"`                   // @description Ключевые слова.
//		Term             string       `json:"term"`                       // @description Условия поиска.
//		Browser          string       `json:"browser"`                    // @description Используемый браузер.
//		Device           string       `json:"device"`                     // @description Тип устройства.
//		OS               string       `json:"os"`                         // @description Операционная система.
//		Template         string       `json:"template"`                   // @description Шаблон страницы.
//		AccountNumber    string       `json:"account_number"`             // @description Номер учетной записи пользователя.
//		BonusSum         float64      `json:"bonus_sum"`                  // @description Сумма бонусов.
//		DeliveryCode     string       `json:"delivery_code"`              // @description Код доставки.
//		DeliveryDate     string       `json:"delivery_date"`              // @description Дата доставки.
//		DeliveryInterval string       `json:"delivery_interval"`          // @description Интервал доставки.
//		Email            string       `json:"email"`                      // @description Электронная почта пользователя.
//		FioCourier       string       `json:"fio_courier"`                // @description ФИО курьера.
//		IsNearestExpress bool         `json:"is_nearest_express"`         // @description Признак экспресс-доставки.
//		Phone            string       `json:"phone"`                      // @description Телефон пользователя.
//		PhoneCourier     string       `json:"phone_courier"`              // @description Телефон курьера.
//		Price            float64      `json:"price"`                      // @description Цена заказа.
//		ShopAddress      string       `json:"shop_address"`               // @description Адрес магазина.
//		PurchaseStatus   string       `json:"purchase_status"`            // @description Статус покупки.
//		ShopSchedule     string       `json:"shop_schedule"`              // @description Расписание магазина.
//		PaymentURL       string       `json:"payment_url"`                // @description URL для оплаты.
//		ProductList      Product_list `json:"product_list"`               // @description Список продуктов в заказе.
//		ShippingCost     float64      `json:"shipping_cost"`              // @description Стоимость доставки.
//		BonusDebet       float64      `json:"bonus_debet"`                // @description Сумма списания бонусов.
//		ShippingAddress  string       `json:"shipping_address"`           // @description Адрес доставки.
//		TrackNumber      string       `json:"track_number"`               // @description Номер для отслеживания.
//		Screen_name      string       `json:"screen_name,omitempty"`      //	@description	Название экрана приложения
//		Screen_classname string       `json:"screen_classname,omitempty"` //	@description	Название класса МП
//		Total_Payment    float64      `json:"total_payment,omitempty"`    //	@description	Общая стоимость заказа с доставкой
//		Is_avatar        bool         `json:"is_avatar,omitempty"`        //	@description	Проверка на аватар
//		//	Device			  string 	`json:"device,omitempty"`	 //	@description	Устройство
//		Device_type string `json:"device_type"` //	@description	Тип устройства
//		//Os          string `json:"os"`          //	@description	ОС устройства
//		//	Browser	          string 	`json:"browser"`			 //	@description	Браузер
//	}
type Notification struct {
	Datetime                   string  `json:"datetime"`
	OrderId                    string  `json:"order_id"`
	ShipmentId                 string  `json:"shipment_id"`
	UserOrderNumber            string  `json:"user_order_number"`
	UserId                     string  `json:"user_id"`
	ParentOrderUpdate          bool    `json:"parent_order_update"`
	Template                   string  `json:"template"`
	Source                     string  `json:"source"`
	Medium                     string  `json:"medium"`
	Campaign                   string  `json:"campaign"`
	Content                    string  `json:"content"`
	Keywords                   string  `json:"keywords"`
	Term                       string  `json:"term"`
	Browser                    string  `json:"browser"`
	Device                     string  `json:"device"`
	Os                         string  `json:"os"`
	AccountNumber              string  `json:"account_number"`
	BonusSum                   int     `json:"bonus_sum"`
	Email                      string  `json:"email"`
	IsNearestExpress           bool    `json:"is_nearest_express"`
	Phone                      string  `json:"phone"`
	Price                      float64 `json:"price"`
	ShippingCost               int     `json:"shipping_cost"`
	TotalPayment               float64 `json:"total_payment"`
	PurchaseStatus             string  `json:"purchase_status"`
	FioCourier                 string  `json:"fio_courier"`
	PhoneCourier               string  `json:"phone_courier"`
	BonusDebet                 int     `json:"bonus_debet"`
	DeliveryCode               string  `json:"delivery_code"`
	DeliveryDate               string  `json:"delivery_date"`
	DeliveryInterval           string  `json:"delivery_interval"`
	TrackNumber                string  `json:"track_number"`
	ShippingAddress            string  `json:"shipping_address"`
	ShopAddress                string  `json:"shop_address"`
	CertificateCustomerMessage string  `json:"certificate_customer_message"`
	CertificateRecipientEmail  string  `json:"certificate_recipient_email"`
	CertificateImageId         string  `json:"certificate_image_id"`
	CertificateImageUrl        string  `json:"certificate_image_url"`
	ProductList                []struct {
		OfferId  string  `json:"offer_id"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	} `json:"product_list"`
	ContactKey string `json:"contact_key"`
	UserEmail  string `json:"user_email"`
}

// Sign_up содержит информацию о регистрации новых пользователей.
// @Description Содержит информацию о регистрации новых пользователей.
// @Tags         sign_up
type Sign_up struct {
	ClientID          string `json:"client_id"`           // @description Уникальный идентификатор пользователя.
	Device_id         string `json:"device_id,omitempty"` //	@description	Идентификатор устройства.
	UserID            string `json:"user_id"`             // @description Уникальный идентификатор зарегистрированного пользователя.
	SessionID         string `json:"session_id"`          // @description Идентификатор сессии пользователя.
	ContactKey        string `json:"contact_key"`         // @description Идентификатор контакта.
	Datetime          string `json:"datetime"`            // @description Дата и время события.
	Source            string `json:"source"`              // @description Источник трафика.
	Medium            string `json:"medium"`              // @description Канал трафика.
	Campaign          string `json:"campaign"`            // @description Название рекламной кампании.
	Content           string `json:"content"`             // @description Содержание кампании.
	Keywords          string `json:"keywords"`            // @description Ключевые слова.
	Term              string `json:"term"`                // @description Условия поиска.
	GeoRegion         string `json:"geo_region"`          // @description Регион.
	GeoCity           string `json:"geo_city"`            // @description Город.
	PagePath          string `json:"page_path"`           // @description Путь к странице.
	PageLocation      string `json:"page_location"`       // @description URL-адрес страницы.
	ScreenName        string `json:"screen_name"`         // @description Название экрана.
	ScreenClassname   string `json:"screen_classname"`    // @description Класс экрана.
	IsExpressDelivery bool   `json:"is_expressdelivery"`  // @description Признак экспресс-доставки.
	PageTitle         string `json:"page_title"`          // @description Заголовок страницы.
	Referrer          string `json:"referrer"`            // @description Реферер.
	IsAuth            bool   `json:"is_auth"`             // @description Признак аутентификации.
	CartID            string `json:"cart_id"`             // @description Идентификатор корзины.
	//Device            string  `json:"device"`              // @description Тип устройства.
	OS string `json:"os"` // @description Операционная система.
	//Browser           string  `json:"browser"`             // @description Используемый браузер.
	EventSource                string  `json:"event_source"`                 // @description Источник события.
	FirstName                  string  `json:"first_name"`                   // @description Имя пользователя.
	SecondName                 string  `json:"second_name"`                  // @description Второе имя пользователя.
	LastName                   string  `json:"last_name"`                    // @description Фамилия пользователя.
	Gender                     string  `json:"gender"`                       // @description Пол пользователя.
	BirthDate                  string  `json:"birth_date"`                   // @description Дата рождения.
	Email                      string  `json:"email"`                        // @description Электронная почта.
	BonusBalance               float64 `json:"bonus_balance"`                // @description Баланс бонусных баллов.
	BonusLevel                 float64 `json:"bonus_level"`                  // @description Уровень бонусной программы.
	CardNumber                 string  `json:"card_number"`                  // @description Номер карты.
	MobilePushOn               bool    `json:"mobile_push_on"`               // @description Признак включения мобильных уведомлений.
	Is_avatar                  bool    `json:"is_avatar,omitempty"`          //	@description	Проверка на аватар
	Device                     string  `json:"device,omitempty"`             //	@description	Устройство
	Device_type                string  `json:"device_type"`                  //	@description	Тип устройства
	Os                         string  `json:"os"`                           //	@description	ОС устройства
	Browser                    string  `json:"browser"`                      //	@description	Браузер
	UserEmail                  string  `json:"user_email"`                   //	@description	email пользолвателя сертификата
	CertificateCustomerMessage string  `json:"certificate_customer_message"` //	@description	сообщение внутри сертификата
	CertificateRecipientEmail  string  `json:"certificate_recipient_email"`  //	@description	email получателя сертифката
	CertificateImageId         string  `json:"certificate_image_id"`         //	@description	id картинки
	CertificateImageUrl        string  `json:"certificate_image_url"`        //	@description	url картинки
}

// Login содержит информацию о входах пользователей в систему.
// @Description Содержит информацию о входах пользователей в систему.
// @Tags         login
type Login struct {
	ClientID          string `json:"client_id"`           // @description Уникальный идентификатор пользователя.
	Device_id         string `json:"device_id,omitempty"` //	@description	Идентификатор устройства.
	UserID            string `json:"user_id"`             // @description Уникальный идентификатор зарегистрированного пользователя.
	SessionID         string `json:"session_id"`          // @description Идентификатор сессии пользователя.
	ContactKey        string `json:"contact_key"`         // @description Идентификатор контакта.
	Datetime          string `json:"datetime"`            // @description Дата и время события.
	Source            string `json:"source"`              // @description Источник трафика.
	Medium            string `json:"medium"`              // @description Канал трафика.
	Campaign          string `json:"campaign"`            // @description Название рекламной кампании.
	Content           string `json:"content"`             // @description Содержание кампании.
	Keywords          string `json:"keywords"`            // @description Ключевые слова.
	Term              string `json:"term"`                // @description Условия поиска.
	GeoRegion         string `json:"geo_region"`          // @description Регион.
	GeoCity           string `json:"geo_city"`            // @description Город.
	PagePath          string `json:"page_path"`           // @description Путь к странице.
	PageLocation      string `json:"page_location"`       // @description URL-адрес страницы.
	ScreenName        string `json:"screen_name"`         // @description Название экрана.
	ScreenClassname   string `json:"screen_classname"`    // @description Класс экрана.
	IsExpressDelivery bool   `json:"is_expressdelivery"`  // @description Признак экспресс-доставки.
	PageTitle         string `json:"page_title"`          // @description Заголовок страницы.
	Referrer          string `json:"referrer"`            // @description Реферер.
	EventSource       string `json:"event_source"`        // @description Источник события.
	IsAuth            bool   `json:"is_auth"`             // @description Признак аутентификации.
	CartID            string `json:"cart_id"`             // @description Идентификатор корзины.
	//Device            string  `json:"device"`              // @description Тип устройства.
	OS string `json:"os"` // @description Операционная система.
	//Browser           string  `json:"browser"`             // @description Используемый браузер.
	FirstName    string  `json:"first_name"`     // @description Имя пользователя.
	SecondName   string  `json:"second_name"`    // @description Второе имя пользователя.
	LastName     string  `json:"last_name"`      // @description Фамилия пользователя.
	Gender       string  `json:"gender"`         // @description Пол пользователя.
	BirthDate    string  `json:"birth_date"`     // @description Дата рождения.
	Email        string  `json:"email"`          // @description Электронная почта.
	BonusBalance float64 `json:"bonus_balance"`  // @description Баланс бонусных баллов.
	BonusLevel   float64 `json:"bonus_level"`    // @description Уровень бонусной программы.
	CardNumber   string  `json:"card_number"`    // @description Номер карты.
	MobilePushOn bool    `json:"mobile_push_on"` // @description Признак включения мобильных уведомлений.
	//Screen_name       string    `json:"screen_name,omitempty"`      //	@description	Название экрана приложения
	//Screen_classname  string    `json:"screen_classname,omitempty"` //	@description	Название класса МП
	Is_avatar   bool   `json:"is_avatar,omitempty"` //	@description	Проверка на аватар
	Device      string `json:"device,omitempty"`    //	@description	Устройство
	Device_type string `json:"device_type"`         //	@description	Тип устройства
	Os          string `json:"os"`                  //	@description	ОС устройства
	Browser     string `json:"browser"`             //	@description	Браузер

}

// @Summary		HttpResponse - это формат HTTP ответа на запрос
// @Description	Ответ бинарный Успех/Неудача, при успехе Код ответа = 0  при неудаче = 1. Текст ошибки содержится в response_text
// @Success
//
//	@Example {
//
// "response_code": 0,
// "response_text": "OK"
// }
//
//	@Example {
//
// "response_code": 1,
// "response_text": "malformed request"
// }
// @Tags			response
type HttpResponse struct {
	ResponseCode int    `json:"response_code"` // @Example 0
	ResponseText string `json:"response_text"` // @Example OK
}

type Add_to_cart struct {
	Client_id               string       `json:"client_id,omitempty"`     //	@description	Идентификатор клиента.
	User_id                 string       `json:"user_id,omitempty"`       //	@description	Идентификатор пользователя.
	Session_id              string       `json:"session_id,omitempty"`    //	@description	Идентификатор сессии.
	Contact_key             string       `json:"contact_key,omitempty"`   //	@description	Идентификатор из Манзана
	Source                  string       `json:"source"`                  //	@description	Источник трафика.
	Medium                  string       `json:"medium"`                  //	@description	Канал трафика.
	Campaign                string       `json:"сampaign"`                //	@description	Кампания.
	Content                 string       `json:"content"`                 //	@description	Контент.
	Datetime                string       `json:"datetime"`                //	@decription Дата и время события
	Keywords                string       `json:"keywords"`                //	@description	Ключевые слова.
	Term                    string       `json:"term"`                    //	@description	Условия поиска.
	Value                   float64      `json:"value"`                   //	@description	Стоимость добавленного товара.
	Currency                string       `json:"currency"`                //	@description	Валюта транзакции.
	Product_id              string       `json:"product_id"`              //	@description	Идентификатор товара.
	Quantity                int          `json:"quantity"`                //	@description	Количество добавленных товаров.
	Cart_id                 string       `json:"cart_id"`                 //	@description	Идентификатор корзины.
	Search_term             string       `json:"search_term"`             //	@description	Поисковый запрос.
	Offer_id                string       `json:"offer_id"`                //	@description	Идентификатор предложения.
	Discount_value          float64      `json:"discount_value"`          //	@description	Сумма скидки на товар.
	Name                    string       `json:"name"`                    //	@description	Название товара.
	Brand                   string       `json:"brand"`                   //	@description	Бренд товара.
	Badge                   string       `json:"badge"`                   //	@description	Значок товара.
	Promo                   string       `json:"promo"`                   //	@description	Промо для товара.
	Rating                  float64      `json:"rating"`                  //	@description	Рейтинг товара.
	Item_list_id            string       `json:"item_list_id"`            //	@description	Идентификатор списка товаров.
	Item_list_name          string       `json:"item_list_name"`          //	@description	Название списка товаров.
	Category_id             int          `json:"category_id"`             //	@description	Идентификатор категории.
	Category_name           string       `json:"category_name"`           //	@description	Название категории.
	Category_hierarchy      string       `json:"category_hierarchy"`      //	@description	Иерархия категорий.
	Product_ids             []string     `json:"product_ids,omitempty"`   //	@description	Идентификаторы товаров в категории.
	Product_list            Product_list `json:"product_list"`            //	@description	Товарный объект с идентификаторами, кол-вом и ценами товаров
	Active_filters          interface{}  `json:"active_filters"`          //	@description	Активные фильтры, применённые пользователем.
	Geo_region              string       `json:"geo_region"`              //	@description	Регион пользователя.
	Geo_city                string       `json:"geo_city"`                //	@description	Город пользователя.
	Page_path               string       `json:"page_path"`               //	@description	Путь страницы.
	Page_location           string       `json:"page_location"`           //	@description	Местоположение страницы.
	Is_expressdelivery      bool         `json:"is_expressdelivery"`      //	@description	Доступна ли экспресс-доставка.
	Page_title              string       `json:"page_title"`              //	@description	Заголовок страницы.
	Referrer                string       `json:"referrer"`                //	@description	Источник перехода на страницу.
	Is_subscriptiondelivery bool         `json:"is_subscriptiondelivery"` //	@description	Доступна ли подписка на доставку.
	Event_source            string       `json:"event_source"`            //	@description	Источник события.
	Is_auth                 bool         `json:"is_auth"`                 //	@description	Авторизован ли пользователь.
	Device_id               string       `json:"device_id,omitempty"`     //	@description	Идентификатор устройства.
	Screen_name             string       `json:"screen_name,omitempty"`   //	@description	Название экрана приложения
	Screen_classname        string       `json:"screen_classname"`        //	@description	Название класса МП
	Is_avatar               bool         `json:"is_avatar,omitempty"`     //	@description	Проверка на аватар
	Device                  string       `json:"device,omitempty"`        //	@description	Устройство
	Device_type             string       `json:"device_type"`             //	@description	Тип устройства
	Os                      string       `json:"os"`                      //	@description	ОС устройства
	Browser                 string       `json:"browser"`                 //	@description	Браузер
}

type View_item struct {
	Client_id               string  `json:"client_id,omitempty"`     //	@description	Идентификатор клиента.
	User_id                 string  `json:"user_id,omitempty"`       //	@description	Идентификатор пользователя.
	Session_id              string  `json:"session_id,omitempty"`    //	@description	Идентификатор сессии.
	Contact_key             string  `json:"contact_key,omitempty"`   //	@description	Идентификатор из Манзана
	Source                  string  `json:"source"`                  //	@description	Источник трафика.
	Medium                  string  `json:"medium"`                  //	@description	Канал трафика.
	Campaign                string  `json:"сampaign"`                //	@description	Кампания.
	Content                 string  `json:"content"`                 //	@description	Контент.
	Datetime                string  `json:"datetime"`                //	@decription Дата и время события
	Keywords                string  `json:"keywords"`                //	@description	Ключевые слова.
	Term                    string  `json:"term"`                    //	@description	Условия поиска.
	Value                   float64 `json:"value"`                   //	@description	Стоимость товара.
	Product_id              string  `json:"product_id"`              //	@description	Идентификатор товара.
	Offer_id                string  `json:"offer_id"`                //	@description	Идентификатор предложения.
	Name                    string  `json:"name"`                    //	@description	Название товара.
	Brand                   string  `json:"brand"`                   //	@description	Бренд товара.
	Currency                string  `json:"currency"`                //	@description	Валюта транзакции.
	Discount_value          float64 `json:"discount_value"`          //	@description	Сумма скидки на товар.
	Badge                   string  `json:"badge"`                   //	@description	Значок товара.
	Promo                   string  `json:"promo"`                   //	@description	Промо для товара.
	Rating                  float64 `json:"rating"`                  //	@description	Рейтинг товара.
	Availability            bool    `json:"availability"`            //	@description	Наличие товара.
	Category_id             int     `json:"category_id"`             //	@description	Идентификатор категории.
	Category_name           string  `json:"category_name"`           //	@description	Название категории.
	Category_hierarchy      string  `json:"category_hierarchy"`      //	@description	Иерархия категорий.
	Product_stm             bool    `json:"product_stm"`             //	@description	Продукт STM.
	Product_food            bool    `json:"product_food"`            //	@description	Продукт питания.
	Pickup_only             bool    `json:"pickup_only"`             //	@description	Только для самовывоза.
	Delivery_date           string  `json:"delivery_date"`           //	@description	Дата доставки.
	Stock_level             int     `json:"stock_level"`             //	@description	Уровень запаса.
	Ratings_count           int     `json:"ratings_count"`           //	@description	Количество оценок.
	Product_spec            string  `json:"product_spec"`            //	@description	Спецификация продукта.
	Product_weight          string  `json:"product_weight"`          //	@description	Вес продукта.
	Product_wear_size       string  `json:"product_wear_size"`       //	@description	Размер продукта.
	Product_color           string  `json:"product_color"`           //	@description	Цвет продукта.
	Product_taste           string  `json:"product_taste"`           //	@description	Вкус продукта.
	Product_pack_count      float64 `json:"product_pack_count"`      //	@description	Количество упаковок.
	Product_pack_type       string  `json:"product_pack_type"`       //	@description	Тип упаковки.
	Product_farma           bool    `json:"product_farma"`           //	@description	Продукт фармацевтики.
	Product_farma_type      string  `json:"product_farma_type"`      //	@description	Тип фармацевтического продукта.
	Pet_type                string  `json:"pet_type"`                //	@description	Тип питомца.
	Pet_age                 string  `json:"pet_age"`                 //	@description	Возраст питомца.
	Pet_size                string  `json:"pet_size"`                //	@description	Размер питомца.
	Geo_region              string  `json:"geo_region"`              //	@description	Регион пользователя.
	Geo_city                string  `json:"geo_city"`                //	@description	Город пользователя.
	Page_path               string  `json:"page_path"`               //	@description	Путь страницы.
	Page_location           string  `json:"page_location"`           //	@description	Местоположение страницы.
	Is_expressdelivery      bool    `json:"is_expressdelivery"`      //	@description	Доступна ли экспресс-доставка.
	Page_title              string  `json:"page_title"`              //	@description	Заголовок страницы.
	Referrer                string  `json:"referrer"`                //	@description	Источник перехода.
	Is_subscriptiondelivery bool    `json:"is_subscriptiondelivery"` //	@description	Доступна ли подписка на доставку.
	Event_source            string  `json:"event_source"`            //	@description	Источник события.
	Is_auth                 bool    `json:"is_auth"`                 //	@description	Авторизован ли пользователь.
	Cart_id                 string  `json:"cart_id"`                 //	@description	Идентификатор корзины.
	Device_id               string  `json:"device_id,omitempty"`     //	@description	Идентификатор устройства.
	Screen_name             string  `json:"screen_name,omitempty"`   //	@description	Название экрана приложения
	Screen_classname        string  `json:"screen_classname"`        //	@description	Название класса МП
	Device                  string  `json:"device,omitempty"`        //	@description	Устройство
	Device_type             string  `json:"device_type"`             //	@description	Тип устройства
	Os                      string  `json:"os"`                      //	@description	ОС устройства
	Browser                 string  `json:"browser"`                 //	@description	Браузер
}
type Begin_checkout struct {
	Client_id   string `json:"client_id"`   //	@description	Идентификатор клиента.
	User_id     string `json:"user_id"`     //	@description	Идентификатор пользователя.
	Session_id  string `json:"session_id"`  //	@description	Идентификатор сессии.
	Contact_key string `json:"contact_key"` //	@description	Ключ контакта.
	Datetime    string `json:"datetime"`    //	@description	Дата и время начала оформления заказа.
	Source      string `json:"source"`      //	@description	Источник трафика.
	Medium      string `json:"medium"`      //	@description	Канал трафика.
	Campaign    string `json:"campaign"`    //	@description	Кампания.
	Content     string `json:"content"`     //	@description	Контент.

	Keywords           string       `json:"keywords"`            //	@description	Ключевые слова.
	Term               string       `json:"term"`                //	@description	Условия поиска.
	Currency           string       `json:"currency"`            //	@description	Валюта транзакции.
	Value              float64      `json:"value"`               //	@description	Общая стоимость товаров.
	Product_list       Product_list `json:"product_list"`        //	@description	Список продуктов.
	Quantity           int          `json:"quantity"`            //	@description	Общее количество товаров.
	Sku_quantity       int          `json:"sku_quantity"`        //	@description	Общее количество SKU.
	Sku_value          float64      `json:"sku_value"`           //	@description	Общая стоимость SKU.
	Order_type         string       `json:"order_type"`          //	@description	Тип заказа.
	Discount_value     float64      `json:"discount_value"`      //	@description	Сумма скидки.
	Geo_region         string       `json:"geo_region"`          //	@description	Регион пользователя.
	Geo_city           string       `json:"geo_city"`            //	@description	Город пользователя.
	Page_path          string       `json:"page_path"`           //	@description	Путь страницы оформления заказа.
	Page_location      string       `json:"page_location"`       //	@description	Местоположение страницы.
	Screen_name        string       `json:"screen_name"`         //	@description	Название экрана.
	Screen_classname   string       `json:"screen_classname"`    //	@description	Класс экрана.
	Is_expressdelivery bool         `json:"is_expressdelivery"`  //	@description	Доступна ли экспресс-доставка.
	Page_title         string       `json:"page_title"`          //	@description	Заголовок страницы.
	Referrer           string       `json:"referrer"`            //	@description	Источник перехода на страницу.
	Event_source       string       `json:"event_source"`        //	@description	Источник события.
	Is_auth            bool         `json:"is_auth"`             //	@description	Авторизован ли пользователь.
	Cart_id            string       `json:"cart_id"`             //	@description	Идентификатор корзины.
	Device_id          string       `json:"device_id,omitempty"` //	@description	Идентификатор устройства.
	Is_avatar          bool         `json:"is_avatar,omitempty"` //	@description	Проверка на аватар
	Device             string       `json:"device,omitempty"`    //	@description	Устройство
	Device_type        string       `json:"device_type"`         //	@description	Тип устройства
	Os                 string       `json:"os"`                  //	@description	ОС устройства
	Browser            string       `json:"browser"`             //	@description	Браузер

}
type Cart_update struct {
	Action                        string       `json:"action"`                         //	@description	Действие (например, добавление или удаление из корзины).
	Button_copy                   string       `json:"button_copy"`                    //	@description	Текст кнопки действия.
	Page_type                     string       `json:"page_type"`                      //	@description	Тип страницы.
	Product_id                    string       `json:"product_id"`                     //	@description	Идентификатор продукта.
	Offer_id                      string       `json:"offer_id"`                       //	@description	Идентификатор предложения.
	Variant_id                    string       `json:"variant_id,omitempty"`           //	@description	Идентификатор варианта продукта.
	Title                         string       `json:"title"`                          //	@description	Название продукта.
	Brand                         string       `json:"brand"`                          //	@description	Бренд продукта.
	Price                         float64      `json:"price"`                          //	@description	Цена продукта.
	Price_local_currency          float64      `json:"price_local_currency,omitempty"` //	@description	Цена продукта в местной валюте.
	Discount_percentage           float64      `json:"discount_percentage"`            //	@description	Процент скидки на продукт.
	Discount_value                float64      `json:"discount_value"`                 //	@description	Сумма скидки на продукт.
	Original_price                float64      `json:"original_price"`                 //	@description	Исходная цена продукта.
	Original_price_local_currency float64      `json:"original_price_local_currency"`  //	@description	Исходная цена продукта в местной валюте.
	Product_list                  Product_list `json:"product_list"`                   //	@description	Список продуктов.
	Product_ids                   []string     `json:"product_ids"`                    //	@description	Идентификаторы продуктов в корзине.
	Total_quantity                int          `json:"total_quantity"`                 //	@description	Общее количество продуктов в корзине.
	Total_price                   float64      `json:"total_price"`                    //	@description	Общая цена всех продуктов в корзине.
	Total_price_without_tax       float64      `json:"total_price_without_tax"`        //	@description	Общая цена продуктов без учёта налогов.
	Total_price_local_currency    float64      `json:"total_price_local_currency"`     //	@description	Общая стоимость продуктов в местной валюте.
	Local_currency                string       `json:"local_currency"`                 //	@description	Местная валюта.
	Tags                          []string     `json:"tags"`                           //	@description	Теги для продуктов.
	Category_1                    string       `json:"category_1"`                     //	@description	Основная категория продукта.
	Category_2                    string       `json:"category_2"`                     //	@description	Вторая категория продукта.
	Category_3                    string       `json:"category_3"`                     //	@description	Третья категория продукта.
	Categories_path               string       `json:"categories_path"`                //	@description	Путь категорий.
	Category_id                   int          `json:"category_id"`                    //	@description	Идентификатор категории.
	Categories_ids                []string     `json:"categories_ids"`                 //	@description	Идентификаторы категорий.
	Variant_list                  Product_list `json:"variant_list"`                   //	@description	Список вариантов продуктов.
	Variant_ids                   []string     `json:"variant_ids"`                    //	@description	Идентификаторы вариантов продуктов.
	Language                      string       `json:"language"`                       //	@description	Язык продукта.
	Location                      string       `json:"location"`                       //	@description	Местоположение.
	Domain                        string       `json:"domain"`                         //	@description	Домен покупки.
	Client_id                     string       `json:"client_id,omitempty"`            //	@description	Идентификатор клиента.
	User_id                       string       `json:"user_id,omitempty"`              //	@description	Идентификатор пользователя.
	Session_id                    string       `json:"session_id,omitempty"`           //	@description	Идентификатор сессии.
	Contact_key                   string       `json:"contact_key,omitempty"`          //	@description	Идентификатор из Манзана
	Device_id                     string       `json:"device_id,omitempty"`            //	@description	Идентификатор устройства.
	Screen_name                   string       `json:"screen_name,omitempty"`          //	@description	Название экрана приложения
	Screen_classname              string       `json:"screen_classname"`               //	@description	Название класса МП
	Is_avatar                     bool         `json:"is_avatar,omitempty"`            //	@description	Проверка на аватар
	Device                        string       `json:"device,omitempty"`               //	@description	Устройство
	Device_type                   string       `json:"device_type"`                    //	@description	Тип устройства
	Os                            string       `json:"os"`                             //	@description	ОС устройства
	Browser                       string       `json:"browser"`                        //	@description	Браузер
}
type View_category struct {
	Client_id               string   `json:"client_id"`               //	@description	Идентификатор клиента.
	User_id                 string   `json:"user_id"`                 //	@description	Идентификатор пользователя.
	Session_id              string   `json:"session_id"`              //	@description	Идентификатор сессии.
	Contact_key             string   `json:"contact_key,omitempty"`   //	@description	Идентификатор из Манзана
	Source                  string   `json:"source"`                  //	@description	Источник трафика.
	Datetime                string   `json:"datetime"`                //	@decription Дата и время события
	Medium                  string   `json:"medium"`                  //	@description	Канал трафика.
	Campaign                string   `json:"сampaign"`                //	@description	Кампания.
	Content                 string   `json:"content"`                 //	@description	Контент.
	Keywords                string   `json:"keywords"`                //	@description	Ключевые слова.
	Term                    string   `json:"term"`                    //	@description	Условия поиска.
	Category_hierarchy      string   `json:"category_hierarchy"`      //	@description	Иерархия категорий.
	Category_id             int      `json:"category_id"`             //	@description	Идентификатор категории.
	Category_name           string   `json:"category_name"`           //	@description	Название категории.
	Page_number             int      `json:"page_number"`             //	@description	Номер страницы.
	Offer_ids               []string `json:"offer_ids"`               //	@description	Идентификаторы продуктов в категории.
	Geo_region              string   `json:"geo_region"`              //	@description	Регион пользователя.
	Geo_city                string   `json:"geo_city"`                //	@description	Город пользователя.
	Page_path               string   `json:"page_path"`               //	@description	Путь страницы.
	Page_location           string   `json:"page_location"`           //	@description	Местоположение страницы.
	Is_expressdelivery      bool     `json:"is_expressdelivery"`      //	@description	Доступна ли экспресс-доставка.
	Page_title              string   `json:"page_title"`              //	@description	Заголовок страницы.
	Referrer                string   `json:"referrer"`                //	@description	Источник перехода.
	Is_subscriptiondelivery bool     `json:"is_subscriptiondelivery"` //	@description	Доступна ли подписка на доставку.
	Event_source            string   `json:"event_source"`            //	@description	Источник события.
	Device_id               string   `json:"device_id,omitempty"`     //	@description	Идентификатор устройства.
	Screen_name             string   `json:"screen_name,omitempty"`   //	@description	Название экрана приложения
	Screen_classname        string   `json:"screen_classname"`        //	@description	Название класса МП
	Cart_id                 string   `json:"cart_id,omitempty"`       //	@description	ID корзины
	Device                  string   `json:"device,omitempty"`        //	@description	Устройство
	Device_type             string   `json:"device_type"`             //	@description	Тип устройства
	Os                      string   `json:"os"`                      //	@description	ОС устройства
	Browser                 string   `json:"browser"`                 //	@description	Браузер
}

type Page_visit struct {
	ClientID          string `json:"client_id"`           // @description Уникальный идентификатор пользователя.
	Device_id         string `json:"device_id,omitempty"` //	@description	Идентификатор устройства.
	UserID            string `json:"user_id"`             // @description Уникальный идентификатор зарегистрированного пользователя.
	SessionID         string `json:"session_id"`          // @description Идентификатор сессии пользователя.
	ContactKey        string `json:"contact_key"`         // @description Идентификатор контакта.
	Datetime          string `json:"datetime"`            // @description Дата и время события.
	Source            string `json:"source"`              // @description Источник трафика.
	Medium            string `json:"medium"`              // @description Канал трафика.
	Campaign          string `json:"campaign"`            // @description Название рекламной кампании.
	Content           string `json:"content"`             // @description Содержание кампании.
	Keywords          string `json:"keywords"`            // @description Ключевые слова.
	Term              string `json:"term"`                // @description Условия поиска.
	GeoRegion         string `json:"geo_region"`          // @description Регион.
	GeoCity           string `json:"geo_city"`            // @description Город.
	PagePath          string `json:"page_path"`           // @description Путь к странице.
	PageLocation      string `json:"page_location"`       // @description URL-адрес страницы.
	ScreenName        string `json:"screen_name"`         // @description Название экрана.
	ScreenClassname   string `json:"screen_classname"`    // @description Класс экрана.
	IsExpressDelivery bool   `json:"is_expressdelivery"`  // @description Признак экспресс-доставки.
	PageTitle         string `json:"page_title"`          // @description Заголовок страницы.
	Referrer          string `json:"referrer"`            // @description Реферер.
	EventSource       string `json:"event_source"`        // @description Источник события.
	IsAuth            bool   `json:"is_auth"`             // @description Признак аутентификации.
	CartID            string `json:"cart_id"`             // @description Идентификатор корзины.
	Is_avatar         bool   `json:"is_avatar,omitempty"` //	@description	Проверка на аватар
	Device            string `json:"device,omitempty"`    //	@description	Устройство
	Device_type       string `json:"device_type"`         //	@description	Тип устройства
	Os                string `json:"os"`                  //	@description	ОС устройства
	Browser           string `json:"browser"`             //	@description	Браузер

}

var requestQueue chan *gin.Context

type Screen_view struct {
	ClientID          string    `json:"client_id"`           // @description Уникальный идентификатор пользователя.
	Device_id         string    `json:"device_id,omitempty"` //	@description	Идентификатор устройства.
	UserID            string    `json:"user_id"`             // @description Уникальный идентификатор зарегистрированного пользователя.
	SessionID         string    `json:"session_id"`          // @description Идентификатор сессии пользователя.
	ContactKey        string    `json:"contact_key"`         // @description Идентификатор контакта.
	Datetime          time.Time `json:"datetime"`            // @description Дата и время события.
	Source            string    `json:"source"`              // @description Источник трафика.
	Medium            string    `json:"medium"`              // @description Канал трафика.
	Campaign          string    `json:"campaign"`            // @description Название рекламной кампании.
	Content           string    `json:"content"`             // @description Содержание кампании.
	Keywords          string    `json:"keywords"`            // @description Ключевые слова.
	Term              string    `json:"term"`                // @description Условия поиска.
	GeoRegion         string    `json:"geo_region"`          // @description Регион.
	GeoCity           string    `json:"geo_city"`            // @description Город.
	PagePath          string    `json:"page_path"`           // @description Путь к странице.
	PageLocation      string    `json:"page_location"`       // @description URL-адрес страницы.
	ScreenName        string    `json:"screen_name"`         // @description Название экрана.
	ScreenClassname   string    `json:"screen_classname"`    // @description Класс экрана.
	IsExpressDelivery bool      `json:"is_expressdelivery"`  // @description Признак экспресс-доставки.
	PageTitle         string    `json:"page_title"`          // @description Заголовок страницы.
	Referrer          string    `json:"referrer"`            // @description Реферер.
	EventSource       string    `json:"event_source"`        // @description Источник события.
	IsAuth            bool      `json:"is_auth"`             // @description Признак аутентификации.
	CartID            string    `json:"cart_id"`             // @description Идентификатор корзины.
	Is_avatar         bool      `json:"is_avatar,omitempty"` //	@description	Проверка на аватар
	Device            string    `json:"device,omitempty"`    //	@description	Устройство
	Device_type       string    `json:"device_type"`         //	@description	Тип устройства
	Os                string    `json:"os"`                  //	@description	ОС устройства
	Browser           string    `json:"browser"`             //	@description	Браузер

}

type EventFileRecord struct {
	Timestamp    time.Time              `json:"timestamp"`
	EventName    string                 `json:"event_name"`
	ResponseCode int                    `json:"response_code"`
	Data         map[string]interface{} `json:"data"`
}
type View_item_list struct {
	ContactKey        string    `json:"contact_key"`         // @description Идентификатор контакта.
	Device_id         string    `json:"device_id,omitempty"` //	@description	Идентификатор устройства.
	Datetime          time.Time `json:"datetime"`            // @description Дата и время события.
	Source            string    `json:"source"`              // @description Источник трафика.
	Medium            string    `json:"medium"`              // @description Канал трафика.
	Campaign          string    `json:"campaign"`            // @description Название рекламной кампании.
	Content           string    `json:"content"`             // @description Содержание кампании.
	Keywords          string    `json:"keywords"`            // @description Ключевые слова.
	Term              string    `json:"term"`                // @description Условия поиска.
	SearchTerm        string    `json:"search_term"`         // @description Условия поиска, указанные в поисковом запросе.
	ItemListName      string    `json:"item_list_name"`      // @description Название списка товаров.
	ItemListID        string    `json:"item_list_id"`        // @description Идентификатор списка товаров.
	CategoryID        string    `json:"category_id"`         // @description Идентификатор категории.
	CategoryName      string    `json:"category_name"`       // @description Название категории.
	CategoryHierarchy string    `json:"category_hierarchy"`  // @description Иерархия категорий.
	OfferIDs          []string  `json:"offer_ids"`           // @description Идентификаторы предложений.
	CartID            string    `json:"cart_id"`             // @description Идентификатор корзины.
	GeoRegion         string    `json:"geo_region"`          // @description Регион.
	GeoCity           string    `json:"geo_city"`            // @description Город.
	PagePath          string    `json:"page_path"`           // @description Путь к странице.
	PageLocation      string    `json:"page_location"`       // @description URL-адрес страницы.
	ScreenName        string    `json:"screen_name"`         // @description Название экрана.
	ScreenClassname   string    `json:"screen_classname"`    // @description Класс экрана.
	IsExpressDelivery bool      `json:"is_expressdelivery"`  // @description Признак экспресс-доставки.
	PageTitle         string    `json:"page_title"`          // @description Заголовок страницы.
	Referrer          string    `json:"referrer"`            // @description Реферер.
	EventSource       string    `json:"event_source"`        // @description Источник события.
	IsAuth            bool      `json:"is_auth"`             // @description Признак аутентификации.
	Is_avatar         bool      `json:"is_avatar,omitempty"` // @description Проверка пользователя
	Device            string    `json:"device,omitempty"`    //	@description	Устройство
	Device_type       string    `json:"device_type"`         //	@description	Тип устройства
	Os                string    `json:"os"`                  //	@description	ОС устройства
	Browser           string    `json:"browser"`             //	@description	Браузер

}
type Digital_reindexer struct {
	UID             int64     `reindex:"uid"`
	ClientID        string    `reindex:"client_id"`
	IP              string    `reindex:"ip"`
	Fingerprint     string    `reindex:"fingerprint"`
	UserAgent       string    `reindex:"user_agent"`
	DeviceID        string    `reindex:"device_id,,pk"`
	TelegramID      string    `reindex:"telegram_id"`
	Location        string    `reindex:"location"`
	SessionID       string    `reindex:"session_id"`
	CreatedAt       time.Time `reindex:"created_at"`
	UpdatedAt       time.Time `reindex:"updated_at"`
	Token           string    `reindex:"token"`
	DeviceLanguage  string    `reindex:"device_language"`
	DevicePlatform  string    `reindex:"device_platform"`
	DeviceTimezone  string    `reindex:"device_timezone"`
	LastAppActivity time.Time `reindex:"last_app_activity"`
	Locale          string    `reindex:"locale"`
}

type LogMessage struct {
	Message string
}
type App_install struct {
	Client_id        string `json:"client_id"`                  //	@description	Идентификатор клиента.
	User_id          string `json:"user_id"`                    //	@description	Идентификатор пользователя.
	Session_id       string `json:"session_id"`                 //	@description	Идентификатор сессии.
	Contact_key      string `json:"contact_key"`                //	@description	Ключ контакта.
	Datetime         string `json:"datetime"`                   //	@description	Дата и время установки приложения.
	Email            string `json:"email"`                      //	@description	Электронная почта пользователя.
	Url              string `json:"url"`                        //	@description	URL страницы, на которой произошло событие.
	Event_source     string `json:"event_source"`               //	@description	Источник события.
	Device_id        string `json:"device_id,omitempty"`        //	@description	Идентификатор устройства.
	Screen_name      string `json:"screen_name,omitempty"`      //	@description	Название экрана приложения
	Screen_classname string `json:"screen_classname,omitempty"` //	@description	Название класса МП
	Device           string `json:"device,omitempty"`           //	@description	Устройство
	Device_type      string `json:"device_type"`                //	@description	Тип устройства
	Os               string `json:"os"`                         //	@description	ОС устройства
	Browser          string `json:"browser"`                    //	@description	Браузер
}

type Purchase struct {
	Client_id                       string       `json:"client_id,omitempty"`               //	@description	Идентификатор клиента.
	Device_id                       string       `json:"device_id,omitempty"`               //	@description	Идентификатор устройства.
	User_id                         string       `json:"user_id,omitempty"`                 //	@description	Идентификатор пользователя.
	Session_id                      string       `json:"session_id,omitempty"`              //	@description	Идентификатор сессии.
	Contact_key                     string       `json:"contact_key,omitempty"`             //	@description	Идентификатор из Манзана
	Purchase_id                     string       `json:"purchase_id"`                       //	@description	Уникальный идентификатор покупки.
	Purchase_status                 string       `json:"purchase_status,omitempty"`         //	@description	Статус покупки (например, завершена, в ожидании).
	Purchase_source_type            string       `json:"purchase_source_type"`              //	@description	Тип источника покупки (например, онлайн, в магазине).
	Purchase_source_name            string       `json:"purchase_source_name,omitempty"`    //	@description	Название источника покупки (например, название магазина или платформы).
	Product_list                    Product_list `json:"product_list"`                      //	@description	Список продуктов, участвующих в покупке.
	Product_ids                     []string     `json:"product_ids"`                       //	@description	Список идентификаторов продуктов, включённых в покупку.
	Total_price                     float64      `json:"total_price"`                       //	@description	Общая стоимость покупки, включая налоги.
	Total_price_without_tax         float64      `json:"total_price_without_tax,omitempty"` //	@description	Общая стоимость покупки без учёта налогов.
	Total_price_local_currency      float64      `json:"total_price_local_currency"`        //	@description	Общая стоимость в местной валюте.
	Local_currency                  string       `json:"local_currency,omitempty"`          //	@description	Местная валюта, использованная для транзакции.
	Total_quantity                  int          `json:"total_quantity"`                    //	@description	Общее количество купленных товаров.
	Payment_type                    string       `json:"payment_type,omitempty"`            //	@description	Способ оплаты (например, кредитная карта, PayPal).
	Shipping_type                   string       `json:"shipping_type,omitempty"`           //	@description	Тип доставки (например, стандартная, экспресс).
	Shipping_company                string       `json:"shipping_company,omitempty"`        //	@description	Компания, ответственная за доставку.
	Shipping_cost                   float64      `json:"shipping_cost,omitempty"`           //	@description	Стоимость доставки.
	Shipping_country                string       `json:"shipping_country,omitempty"`        //	@description	Страна, в которую осуществляется доставка.
	Shipping_city                   string       `json:"shipping_city,omitempty"`           //	@description	Город, в который осуществляется доставка.
	Tax_percentage                  float64      `json:"tax_percentage,omitempty"`          //	@description	Процент налога, применённого к покупке.
	Tax_value                       float64      `json:"tax_value,omitempty"`               //	@description	Общая сумма налога, применённого к покупке.
	Coupon_code                     string       `json:"coupon_code,omitempty"`             //	@description	Промокод, использованный при покупке, если он был.
	Voucher_percentage              float64      `json:"voucher_percentage,omitempty"`      //	@description	Процент скидки, применённой с помощью промокода.
	Coupon_value                    float64      `json:"coupon_value,omitempty"`            //	@description	Общая сумма скидки, применённой с помощью промокода.
	Variant_list                    []string     `json:"variant_list,omitempty"`            //	@description	Список вариантов продуктов, включённых в покупку.
	Variant_ids                     []string     `json:"variant_ids,omitempty"`             //	@description	Список идентификаторов вариантов продуктов.
	Language                        string       `json:"language,omitempty"`                //	@description	Язык, использованный при оформлении покупки.
	Location                        string       `json:"location"`                          //	@description	Географическое местоположение, где была совершена покупка.
	Is_waiting_for_delivery_to_shop bool         `json:"is_waiting_for_delivery_to_shop"`   //	@description	Ждет доставки в магазин
	Was_delivered_to_shop           bool         `json:"was_delivered_to_shop"`             //	@description	Был доставлен в магазин
	Domain                          string       `json:"domain"`                            //	@description	Домен или вебсайт, с которого была совершена покупка.
	Screen_name                     string       `json:"screen_name,omitempty"`             //	@description	Название экрана приложения
	Screen_classname                string       `json:"screen_classname"`                  //	@description	Название класса МП
	Currency                        string       `json:"currency"`                          // @description Валюта.
	Value                           float64      `json:"value,omitempty"`                   // @description Стоимость.
	TransactionDate                 string       `json:"transaction_date"`                  // @description Дата транзакции.
	SinglePackage                   bool         `json:"single_package,omitempty"`          // @description Признак единого пакета.
	RecipientType                   string       `json:"recipient_type,omitempty"`          // @description Тип получателя.
	BonusesSpentValue               float64      `json:"bonuses_spent_value,omitempty"`     // @description Потраченные бонусы.
	BonusesAddedValue               float64      `json:"bonuses_added_value,omitempty"`     // @description Начисленные бонусы.
	TotalWeight                     float64      `json:"total_weight,omitempty"`            // @description Общий вес заказа.
	Is_avatar                       bool         `json:"is_avatar,omitempty"`               //	@description	Проверка на аватар
	Device                          string       `json:"device,omitempty"`                  //	@description	Устройство
	Device_type                     string       `json:"device_type"`                       //	@description	Тип устройства
	Os                              string       `json:"os"`                                //	@description	ОС устройства
	Browser                         string       `json:"browser"`                           //	@description	Браузер
}
type Purchase_items struct {
	Purchase_id                   string   `json:"purchase_id"`                             //	@description	Уникальный идентификатор покупки.
	Purchase_status               string   `json:"purchase_status"`                         //	@description	Статус покупки.
	Purchase_source_type          string   `json:"purchase_source_type"`                    //	@description	Тип источника покупки.
	Purchase_source_name          string   `json:"purchase_source_name,omitempty"`          //	@description	Название источника покупки.
	Product_id                    string   `json:"product_id"`                              //	@description	Уникальный идентификатор продукта.
	Variant_id                    string   `json:"variant_id"`                              //	@description	Идентификатор варианта продукта.
	Title                         string   `json:"title"`                                   //	@description	Название продукта.
	Brand                         string   `json:"brand"`                                   //	@description	Бренд продукта.
	Price                         float64  `json:"price"`                                   //	@description	Цена продукта.
	Price_local_currency          float64  `json:"price_local_currency"`                    //	@description	Цена продукта в местной валюте.
	Discount_percentage           float64  `json:"discount_percentage,omitempty"`           //	@description	Процент скидки на продукт.
	Discount_value                float64  `json:"discount_value,omitempty"`                //	@description	Сумма скидки на продукт.
	Original_price                float64  `json:"original_price,omitempty"`                //	@description	Исходная цена продукта до скидки.
	Original_price_local_currency float64  `json:"original_price_local_currency,omitempty"` //	@description	Исходная цена продукта в местной валюте.
	Quantity                      int      `json:"quantity"`                                //	@description	Количество продуктов.
	Total_price                   float64  `json:"total_price"`                             //	@description	Общая стоимость продуктов.
	Total_price_without_tax       float64  `json:"total_price_without_tax"`                 //	@description	Общая стоимость продуктов без учёта налогов.
	Total_price_local_currency    float64  `json:"total_price_local_currency"`              //	@description	Общая стоимость продуктов в местной валюте.
	Local_currency                string   `json:"local_currency"`                          //	@description	Местная валюта.
	Tags                          []string `json:"tags,omitempty"`                          //	@description	Теги для продукта.
	Category_1                    string   `json:"category_1"`                              //	@description	Основная категория продукта.
	Category_2                    string   `json:"category_2"`                              //	@description	Вторая категория продукта.
	Category_3                    string   `json:"category_3"`                              //	@description	Третья категория продукта.
	Categories_path               string   `json:"categories_path"`                         //	@description	Путь по категориям.
	Category_id                   int      `json:"category_id"`                             //	@description	Идентификатор категории.
	Categories_ids                string   `json:"categories_ids,omitempty"`                //	@description	Идентификаторы категорий.
	Language                      string   `json:"language"`                                //	@description	Язык продукта.
	Location                      string   `json:"location"`                                //	@description	Местоположение покупки.
	Domain                        string   `json:"domain"`                                  //	@description	Домен покупки.
	Client_id                     string   `json:"client_id,omitempty"`                     //	@description	Идентификатор клиента.
	User_id                       string   `json:"user_id,omitempty"`                       //	@description	Идентификатор пользователя.
	Session_id                    string   `json:"session_id,omitempty"`                    //	@description	Идентификатор сессии.
	Contact_key                   string   `json:"contact_key,omitempty"`                   //	@description	Идентификатор из Манзана
	Device_id                     string   `json:"device_id,omitempty"`                     //	@description	Идентификатор устройства.
	Screen_name                   string   `json:"screen_name,omitempty"`                   //	@description	Название экрана приложения
	Screen_classname              string   `json:"screen_classname"`                        //	@description	Название класса МП
	Device                        string   `json:"device,omitempty"`                        //	@description	Устройство
	Device_type                   string   `json:"device_type"`                             //	@description	Тип устройства
	Os                            string   `json:"os"`                                      //	@description	ОС устройства
	Browser                       string   `json:"browser"`                                 //	@description	Браузер
}

func StructGuesser(c *gin.Context, event_files_log EventsFileDic, deviceMap map[string]Digital_reindexer, reindexDB *reindexer.Reindexer) (interface{}, string, error) {

	var L LogMessage
	start_time := time.Now()
	//c.Copy()
	switch c.Param("name") {

	case "view_item_list":

		var json_input View_item_list

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {
			//input, _ := io.ReadAll(c.Request.Body)
			//fmt.Println("incoming data: ", string(input))
			//c.JSON(http.StatusBadRequest,

			//	struct {
			//		ResponseCode int    `json:"response_code"`
			//		ResponseText string `json:"response_text"`
			//	}{
			//		ResponseCode: 1,
			//		ResponseText: "malformed request",
			//	})
			//end_time := time.Now().Sub(start_time)
			//L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			//L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input View_item_list, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input.Client_id)

			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			//c.JSON(http.StatusOK, struct {
			//	ResponseCode int    `json:"response_code"`
			//	ResponseText string `json:"response_text"`
			//}{
			//	ResponseCode: 0,
			//	ResponseText: "OK",
			//})
			//input, _ := io.ReadAll(c.Request.Body)
			//fmt.Println("incoming data: ", string(input))
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("result", true)
			c.Set("payload", json_input)
			go func(json_input View_item_list, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "view_item_list", nil
		}

	case "page_visit":

		var json_input Page_visit

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {
			//input, _ := io.ReadAll(c.Request.Body)
			//fmt.Println("incoming data: ", string(input))
			//c.JSON(http.StatusBadRequest,

			//	struct {
			//		ResponseCode int    `json:"response_code"`
			//		ResponseText string `json:"response_text"`
			//	}{
			//		ResponseCode: 1,
			//		ResponseText: "malformed request",
			//	})
			//end_time := time.Now().Sub(start_time)
			//L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			//L.Logger()
			////LogInFile(FileStorage[filename], L.Message)
			//fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Page_visit, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input.Client_id)

			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			//c.JSON(http.StatusOK, struct {
			//	ResponseCode int    `json:"response_code"`
			//	ResponseText string `json:"response_text"`
			//}{
			//	ResponseCode: 0,
			//	ResponseText: "OK",
			//})
			//input, _ := io.ReadAll(c.Request.Body)
			//fmt.Println("incoming data: ", string(input))
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("result", true)
			c.Set("payload", json_input)
			go func(json_input Page_visit, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "page_visit", nil
		}
	case "screen_view":

		var json_input Screen_view

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {
			fmt.Println(err)
			//c.JSON(http.StatusBadRequest,
			//input, _ := io.ReadAll(c.Request.Body)
			//fmt.Println("incoming data: ", string(input))
			//	struct {
			//		ResponseCode int    `json:"response_code"`
			//		ResponseText string `json:"response_text"`
			//	}{
			//		ResponseCode: 1,
			//		ResponseText: "malformed request",
			//	})
			//end_time := time.Now().Sub(start_time)
			//L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			//L.Logger()
			////LogInFile(FileStorage[filename], L.Message)
			//fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Screen_view, c *gin.Context, deviceMap map[string]Digital_reindexer, reindexDB *reindexer.Reindexer) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)

				_, ok := deviceMap[json_input.Device_id]
				if ok {

					if json_input.ContactKey != "" {
						deviceMap[json_input.Device_id] = Digital_reindexer{
							UID:             deviceMap[json_input.Device_id].UID,
							ClientID:        json_input.ContactKey,
							IP:              deviceMap[json_input.Device_id].IP,
							Fingerprint:     "",
							UserAgent:       json_input.Device_type,
							DeviceID:        deviceMap[json_input.Device_id].DeviceID,
							TelegramID:      "",
							Location:        "",
							SessionID:       "",
							CreatedAt:       deviceMap[json_input.Device_id].CreatedAt,
							UpdatedAt:       time.Now(),
							Token:           deviceMap[json_input.Device_id].Token,
							DeviceLanguage:  "RU",
							DevicePlatform:  json_input.Os,
							DeviceTimezone:  deviceMap[json_input.Device_id].DeviceTimezone,
							LastAppActivity: time.Now(),
							Locale:          "RU",
						}

					}

					reindexDB.Upsert("digital", Digital_reindexer{
						UID:             deviceMap[json_input.Device_id].UID,
						ClientID:        json_input.ContactKey,
						IP:              deviceMap[json_input.Device_id].IP,
						Fingerprint:     "",
						UserAgent:       json_input.Device_type,
						DeviceID:        json_input.Device_id,
						TelegramID:      "",
						Location:        "",
						SessionID:       "",
						CreatedAt:       deviceMap[json_input.Device_id].CreatedAt,
						UpdatedAt:       time.Now(),
						Token:           deviceMap[json_input.Device_id].Token,
						DeviceLanguage:  "RU",
						DevicePlatform:  json_input.Os,
						DeviceTimezone:  deviceMap[json_input.Device_id].DeviceTimezone,
						LastAppActivity: time.Now(),
						Locale:          "RU",
					})

					//reindexDB.Close()
				}

				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c, deviceMap, reindexDB)
			requestQueue <- c

			return nil, "", err
		} else {

			c_k, err := CheckContact_Key(json_input)

			if err != nil {

				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("result", true)
			c.Set("payload", json_input)

			go func(json_input Screen_view, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)

				_, ok := deviceMap[json_input.Device_id]
				if ok {

					if json_input.ContactKey != "" {
						deviceMap[json_input.Device_id] = Digital_reindexer{
							UID:             deviceMap[json_input.Device_id].UID,
							ClientID:        json_input.ContactKey,
							IP:              deviceMap[json_input.Device_id].IP,
							Fingerprint:     "",
							UserAgent:       json_input.Device_type,
							DeviceID:        deviceMap[json_input.Device_id].DeviceID,
							TelegramID:      "",
							Location:        "",
							SessionID:       "",
							CreatedAt:       deviceMap[json_input.Device_id].CreatedAt,
							UpdatedAt:       time.Now(),
							Token:           deviceMap[json_input.Device_id].Token,
							DeviceLanguage:  "RU",
							DevicePlatform:  json_input.Os,
							DeviceTimezone:  deviceMap[json_input.Device_id].DeviceTimezone,
							LastAppActivity: time.Now(),
							Locale:          "RU",
						}

					}

					reindexDB.Upsert("digital", Digital_reindexer{
						UID:             deviceMap[json_input.Device_id].UID,
						ClientID:        json_input.ContactKey,
						IP:              deviceMap[json_input.Device_id].IP,
						Fingerprint:     "",
						UserAgent:       json_input.Device_type,
						DeviceID:        deviceMap[json_input.Device_id].DeviceID,
						TelegramID:      "",
						Location:        "",
						SessionID:       "",
						CreatedAt:       deviceMap[json_input.Device_id].CreatedAt,
						UpdatedAt:       time.Now(),
						Token:           deviceMap[json_input.Device_id].Token,
						DeviceLanguage:  "RU",
						DevicePlatform:  json_input.Os,
						DeviceTimezone:  deviceMap[json_input.Device_id].DeviceTimezone,
						LastAppActivity: time.Now(),
						Locale:          "RU",
					})

					//reindexDB.Close()
				}
				//reindexDB.Upsert("digital", Digital_reindexer{
				//
				//	ClientID: json_input.ContactKey,
				//
				//	UserAgent:  json_input.Device_type,
				//	DeviceID:   json_input.Device_id,
				//	TelegramID: "",
				//	Location:   "",
				//	SessionID:  "",
				//
				//	UpdatedAt: time.Now(),
				//
				//	DeviceLanguage: "RU",
				//	DevicePlatform: json_input.Os,
				//
				//	LastAppActivity: time.Now(),
				//	Locale:          "RU",
				//})
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)

			requestQueue <- c

			return json_input, "screen_view", nil
		}

	case "view_category":

		var json_input View_category

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			fmt.Println(err)
			//c.JSON(http.StatusBadRequest,
			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))
			//	struct {
			//		ResponseCode int    `json:"response_code"`
			//		ResponseText string `json:"response_text"`
			//	}{
			//		ResponseCode: 1,
			//		ResponseText: "malformed request",
			//	})
			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input View_category, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input.Client_id)

			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			//c.JSON(http.StatusOK, struct {
			//	ResponseCode int    `json:"response_code"`
			//	ResponseText string `json:"response_text"`
			//}{
			//	ResponseCode: 0,
			//	ResponseText: "OK",
			//})
			//input, _ := io.ReadAll(v.Request.Body)
			//fmt.Println("incoming data: ", string(input))
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("result", true)
			c.Set("payload", json_input)
			go func(json_input View_category, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "view_category", nil
		}
	case "purchase":

		var json_input Purchase

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {
			//input, _ := io.ReadAll(c.Request.Body)
			//fmt.Println("incoming data: ", string(input))
			//end_time := time.Now().Sub(start_time)
			//L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			//L.Logger()
			////LogInFile(FileStorage[filename], L.Message)
			//fmt.Println(L)

			c.Set("error", err.Error())
			go func(json_input Purchase, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("result", true)
			c.Set("payload", json_input)

			//fmt.Println("Событие Purchase исходный json: ")
			go func(json_input Purchase, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)

			requestQueue <- c
			return json_input, "purchase", nil
		}
	case "purchase_items":

		var json_input Purchase_items

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			//input, _ := io.ReadAll(c.Request.Body)
			//fmt.Println("incoming data: ", string(input))
			//
			//end_time := time.Now().Sub(start_time)
			//L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			//L.Logger()
			////LogInFile(FileStorage[filename], L.Message)
			//fmt.Println(L)
			c.Set("error", err.Error())

			go func(json_input Purchase_items, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("result", true)
			c.Set("payload", json_input)
			go func(json_input Purchase_items, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "purchase_items", nil
		}
	case "cart_update":

		var json_input Cart_update

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)
			c.Set("error", err.Error())

			go func(json_input Cart_update, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Cart_update, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "cart_update", nil
		}
	case "view_item":

		var json_input View_item

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {
			//c.JSON(http.StatusBadRequest,
			//	struct {
			//		ResponseCode int    `json:"response_code"`
			//		ResponseText string `json:"response_text"`
			//	}{
			//		ResponseCode: 1,
			//		ResponseText: "malformed request",
			//	})

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))
			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input View_item, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)

			_, err = gin.DefaultWriter.Write([]byte(log_message))
			if err != nil {
				fmt.Println(err)

			}
			//c.JSON(http.StatusOK, struct {
			//	ResponseCode int    `json:"response_code"`
			//	ResponseText string `json:"response_text"`
			//}{
			//	ResponseCode: 0,
			//	ResponseText: "OK",
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			//			c.Next()
			go func(json_input View_item, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "view_item", nil
		}
	case "add_to_cart":

		var json_input Add_to_cart

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {
			//c.JSON(http.StatusBadRequest,
			//	struct {
			//		ResponseCode int    `json:"response_code"`
			//		ResponseText string `json:"response_text"`
			//	}{
			//		ResponseCode: 1,
			//		ResponseText: "malformed request",
			//	})
			//end_time := time.Now().Sub(start_time)

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Add_to_cart, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)

			_, err = gin.DefaultWriter.Write([]byte(log_message))
			if err != nil {
				fmt.Println(err)

			}
			//c.JSON(http.StatusOK, struct {
			//	ResponseCode int    `json:"response_code"`
			//	ResponseText string `json:"response_text"`
			//}{
			//	ResponseCode: 0,
			//	ResponseText: "OK",
			//})
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Add_to_cart, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "add_to_cart", nil
		}
	case "remove_from_cart":

		var json_input Remove_from_cart

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {
			//c.JSON(http.StatusBadRequest,
			//	struct {
			//		ResponseCode int    `json:"response_code"`
			//		ResponseText string `json:"response_text"`
			//	}{
			//		ResponseCode: 1,
			//		ResponseText: "malformed request",
			//	})
			//end_time := time.Now().Sub(start_time)
			//L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Remove_from_cart, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)

			_, err = gin.DefaultWriter.Write([]byte(log_message))
			if err != nil {
				fmt.Println(err)

			}
			//c.JSON(http.StatusOK, struct {
			//	ResponseCode int    `json:"response_code"`
			//	ResponseText string `json:"response_text"`
			//}{
			//	ResponseCode: 0,
			//	ResponseText: "OK",
			//})
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Remove_from_cart, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "remove_from_cart", nil
		}
	case "add_to_wishlist":

		var json_input Add_to_Wishlist

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {
			//c.JSON(http.StatusBadRequest,
			//	struct {
			//		ResponseCode int    `json:"response_code"`
			//		ResponseText string `json:"response_text"`
			//	}{
			//		ResponseCode: 1,
			//		ResponseText: "malformed request",
			//	})
			//end_time := time.Now().Sub(start_time)
			//L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			//L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			//fmt.Println(L)

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Add_to_Wishlist, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)

			_, err = gin.DefaultWriter.Write([]byte(log_message))
			if err != nil {
				fmt.Println(err)

			}
			//c.JSON(http.StatusOK, struct {
			//	ResponseCode int    `json:"response_code"`
			//	ResponseText string `json:"response_text"`
			//}{
			//	ResponseCode: 0,
			//	ResponseText: "OK",
			//})
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Add_to_Wishlist, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "add_to_wishlist", nil
		}
	case "view_cart":

		var json_input View_Cart

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {
			//c.JSON(http.StatusBadRequest,
			//	struct {
			//		ResponseCode int    `json:"response_code"`
			//		ResponseText string `json:"response_text"`
			//	}{
			//		ResponseCode: 1,
			//		ResponseText: "malformed request",
			//	})
			//end_time := time.Now().Sub(start_time)
			//L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input View_Cart, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)

			_, err = gin.DefaultWriter.Write([]byte(log_message))
			if err != nil {
				fmt.Println(err)

			}
			//c.JSON(http.StatusOK, struct {
			//	ResponseCode int    `json:"response_code"`
			//	ResponseText string `json:"response_text"`
			//}{
			//	ResponseCode: 0,
			//	ResponseText: "OK",
			//})
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			//c.Set("payload", json_input)
			//c.Next()
			go func(json_input View_Cart, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "view_cart", nil
		}
	case "begin_checkout":

		var json_input Begin_checkout

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {
			//c.JSON(http.StatusBadRequest,
			//	struct {
			//		ResponseCode int    `json:"response_code"`
			//		ResponseText string `json:"response_text"`
			//	}{
			//		ResponseCode: 1,
			//		ResponseText: "malformed request",
			//	})
			//end_time := time.Now().Sub(start_time)
			//L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			//L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			//fmt.Println(L)

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Begin_checkout, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			//c.JSON(http.StatusOK, struct {
			//	ResponseCode int    `json:"response_code"`
			//	ResponseText string `json:"response_text"`
			//}{
			//	ResponseCode: 0,
			//	ResponseText: "OK",
			//})
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Begin_checkout, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "begin_checkout", nil
		}
	case "add_contact_info":

		var json_input Add_contact_info

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Add_contact_info, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Add_contact_info, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "add_contact_info", nil
		}
	case "add_payment_info":

		var json_input Add_payment_info

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Add_payment_info, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			go func(json_input Add_payment_info, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "add_payment_info", nil
		}
	case "finish_checkout":

		var json_input Finish_checkout

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Finish_checkout, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			go func(json_input Finish_checkout, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "finish_checkout", nil
		}
	case "double_opt_in":

		var json_input Double_opt_in

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Double_opt_in, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Double_opt_in, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "double_opt_in", nil
		}
	case "opt_in_confirmed":

		var json_input Opt_in_confirmed

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)
			c.Set("error", err.Error())

			go func(json_input Opt_in_confirmed, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Opt_in_confirmed, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "opt_in_confirmed", nil
		}
	case "app_install":

		var json_input App_install

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input App_install, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input App_install, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "app_install", nil
		}
	case "purchase_item":

		var json_input Purchase_item

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Purchase_item, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Purchase_item, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "purchase_item", nil
		}
	case "add_shipping_info":

		var json_input Add_shipping_info

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Add_shipping_info, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Add_shipping_info, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "add_shiping_info", nil
		}

	case "express_delivery":

		var json_input Express_delivery
		//body_byte, _ := c.GetRawData()
		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			//fmt.Println(L)

			//fmt.Println("получены предварительно данные ", string(body_byte))
			c.Set("error", err.Error())

			go func(json_input Express_delivery, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {
				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Express_delivery, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "express_delivery", nil
		}
	case "comment":

		var json_input Comment

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Comment, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {

				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Comment, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "comment", nil
		}
	case "search":

		var json_input Search

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)
			c.Set("error", err.Error())

			go func(json_input Search, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {

				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Search, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "search", nil
		}
	case "notification":

		var json_input Notification

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())
			go func(json_input Notification, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {

				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			//c.Next()
			go func(json_input Notification, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "notification", nil
		}
	case "sign_up":

		var json_input Sign_up

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Sign_up, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {

				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Sign_up, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "sign_up", nil
		}
	case "login":

		var json_input Login

		if err := c.ShouldBindBodyWith(&json_input, binding.JSON); err != nil {

			input, _ := io.ReadAll(c.Request.Body)
			fmt.Println("incoming data: ", string(input))

			end_time := time.Now().Sub(start_time)
			L.Message = "Account  couldn't be activated due to error: " + err.Error() + "in " + strconv.Itoa(int(end_time.Microseconds())) + " µs"
			L.Logger()
			//LogInFile(FileStorage[filename], L.Message)
			fmt.Println(L)

			c.Set("error", err.Error())

			go func(json_input Login, c *gin.Context) {
				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)

				eventFileRecord := EventFileRecord{}
				eventFileRecord.Set(input_data, c.Param("name"), 400)

				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "error",
						"event_name":   c.Param("name"),
					},
					"error",
					string(json_byte))

			}(json_input, c)
			requestQueue <- c

			return nil, "", err
		} else {
			//byteBody, _ := ioutil.ReadAll(c.Request.Body)
			//log_message := fmt.Sprintf("event %s was succesfully proccessed", json_input)
			//
			//_, err = gin.DefaultWriter.Write([]byte(log_message))
			//if err != nil {
			//	fmt.Println(err)
			//
			//}
			c_k, err := CheckContact_Key(json_input)

			if err != nil {

				input, _ := io.ReadAll(c.Request.Body)
				fmt.Println("incoming data: ", string(input))
				fmt.Println("Ошибка: ", err)

			}

			c.Set("contact_key", c_k)
			c.Set("payload", json_input)
			c.Next()
			go func(json_input Login, c *gin.Context) {

				input_data := make(map[string]interface{})
				json_byte, _ := json.Marshal(json_input)
				//val, _ := c.Get("payload")

				json.Unmarshal(json_byte, &input_data)
				//fmt.Println("PAYLOAD GO:", string(json_byte))

				eventFileRecord := EventFileRecord{}

				fmt.Println("Успешное получение события " + c.Param("name"))
				eventFileRecord.Set(input_data, c.Param("name"), 200)
				fmt.Println("EFR :", fmt.Sprintf("%+v", eventFileRecord))
				event_files_log.UpdateFileWithRecord(c.Param("name"), eventFileRecord)
				RecordLoki(
					map[string]string{
						"events":       "events",
						"event_status": "success",
						"event_name":   c.Param("name"),
					},
					"info",
					string(json_byte))

				return
			}(json_input, c)
			requestQueue <- c

			return json_input, "login", nil
		}
	}
	c.Set("result", false)
	c.Next()
	return nil, "", errors.New("no event name was passed as parameter")
}
func (L *LogMessage) Logger() {

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	L.Message = "[" + timestamp + "] " + L.Message + "\n"

}

type EventsFileDic map[string]string

func createFileWithAbsolutePath(dirPath, filename string) error {

	//if err := os.MkdirAll(dirPath, 0755); err != nil {
	//	return err
	//}

	filePath := filepath.Join(dirPath, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}
func CreateFileLog(event_name string, EFD EventsFileDic) error {

	err := createFileWithAbsolutePath(events_log_folder, EFD[event_name])
	fmt.Println("Создаем файл ", EFD[event_name])
	if err != nil {

		return err

	}

	return nil

}

var (
	events_log_folder = "/home/fla/blogs"

	event_name_map = map[string]string{

		"view_category":     "Просмотр категории",
		"purchase":          "Покупка",
		"purchase_items":    "Покупка с указанием товаров",
		"cart_update":       "Обновление корзины",
		"view_item":         "Просмотр товарной позиции",
		"add_to_cart":       "Добавление в корзину",
		"select_item":       "Клик на рекомендательный блок",
		"view_item_list":    "Просмотр рекомендательного блока",
		"remove_from_cart":  "Удаление товара из корзины",
		"add_to_wishlist":   "Добавление в избранное",
		"view_cart":         "Просмотр корзины",
		"begin_checkout":    "Начало оформления заказа",
		"add_contact_info":  "Добавление контактных данных",
		"add_payment_info":  "Добавление данных платежа",
		"finish_checkout":   "Завершение оформления заказа",
		"double_opt_in":     "При появлении в базе нового полльзователя",
		"opt_in_confirmed":  "При подтверждении email пользователем",
		"app_install":       "Первый запуск приложения",
		"purchase_item":     "Оформление одной единицы товара",
		"add_shipping_info": "Добавление информации о доставке",
		"express_delivery":  "Экспресс доставка",
		"comment":           "После публикации отзыва на товар",
		"search":            "Поисковой запрос",
		"notification":      "Уведомление",
		"comment_reaction":  "Реакция на комментарий",
		"sign_up":           "При успешной регистрации",
		"login":             "Логин на сайте/ МП",
		"page_visit":        "Посещение страницы сайта",
		"screen_view":       "Посещение экрана мобильного приложениия",
		"unsubscribe_email": "Отписка от email рассылки",
		"unsubscribe_sms":   "Отписка от sms рассылки",
		"unsubscribe_push":  "Отписка от PUSH уведомлений",
	}

	VK_CONFIG = TokenResponseVK{}
)

func (EFD EventsFileDic) UpdateFileWithRecord(event_name string, eventFileRecord EventFileRecord) error {
	f, err := os.OpenFile(filepath.Join("../blogs", EFD[event_name]),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		return err
	}
	defer f.Close()

	data_bytes, _ := json.Marshal(eventFileRecord)

	writer := bufio.NewWriter(f)

	if _, err := writer.WriteString(string(data_bytes) + "\n"); err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
func (E *EventFileRecord) Set(data map[string]interface{}, event_name string, responseCode int) {
	E.Timestamp = time.Now()
	E.EventName = event_name
	E.ResponseCode = responseCode
	E.Data = data
}

func CheckContact_Key(body interface{}) (string, error) {
	var result map[string]any

	input, _ := json.Marshal(body)

	//fmt.Println("Наличие идентификатора ", string(input))

	json.Unmarshal(input, &result)

	_, ok := result["contact_key"].(interface{})
	if ok {

		val, ok2 := result["contact_key"].(interface{}).(string)
		if ok2 {

			return val, nil

		} else {
			return "", errors.New("contact_key is missing")

		}

	} else {
		return "", errors.New("contact_key is missing")
	}

}
