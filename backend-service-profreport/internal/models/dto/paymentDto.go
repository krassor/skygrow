package dto

// PayRequestDto структура для получения уведомлений от CloudPayments
type PayRequestDto struct {
	// Обязательные параметры
	TransactionId   int64   `json:"TransactionId"`   // Номер транзакции в системе
	Amount          float64 `json:"Amount"`          // Сумма оплаты из параметров платежа
	Currency        string  `json:"Currency"`        // Валюта: RUB/USD/EUR/GBP
	PaymentAmount   string  `json:"PaymentAmount"`   // Сумма списания
	PaymentCurrency string  `json:"PaymentCurrency"` // Валюта списания
	DateTime        string  `json:"DateTime"`        // Дата/время создания платежа UTC (yyyy-MM-dd HH:mm:ss)
	CardFirstSix    string  `json:"CardFirstSix"`    // Первые 6 цифр номера карты
	CardLastFour    string  `json:"CardLastFour"`    // Последние 4 цифры номера карты
	CardType        string  `json:"CardType"`        // Платежная система: Visa, Mastercard, Maestro, МИР
	CardExpDate     string  `json:"CardExpDate"`     // Срок действия карты MM/YY
	TestMode        int     `json:"TestMode"`        // Признак тестового режима (1 или 0)
	Status          string  `json:"Status"`          // Статус: Completed/Authorized
	OperationType   string  `json:"OperationType"`   // Тип операции: Payment/CardPayout
	GatewayName     string  `json:"GatewayName"`     // Идентификатор банка-эквайера

	// Необязательные параметры
	CardId                                string                 `json:"CardId,omitempty"`                                // Уникальный идентификатор карты
	InvoiceId                             string                 `json:"InvoiceId,omitempty"`                             // Номер заказа
	AccountId                             string                 `json:"AccountId,omitempty"`                             // Идентификатор пользователя
	SubscriptionId                        string                 `json:"SubscriptionId,omitempty"`                        // Идентификатор подписки
	Name                                  string                 `json:"Name,omitempty"`                                  // Имя держателя карты
	Email                                 string                 `json:"Email,omitempty"`                                 // E-mail адрес плательщика
	IpAddress                             string                 `json:"IpAddress,omitempty"`                             // IP-адрес плательщика
	IpCountry                             string                 `json:"IpCountry,omitempty"`                             // Двухбуквенный код страны ISO3166-1
	IpCity                                string                 `json:"IpCity,omitempty"`                                // Город нахождения плательщика
	IpRegion                              string                 `json:"IpRegion,omitempty"`                              // Регион нахождения плательщика
	IpDistrict                            string                 `json:"IpDistrict,omitempty"`                            // Округ нахождения плательщика
	IpLatitude                            string                 `json:"IpLatitude,omitempty"`                            // Широта нахождения плательщика
	IpLongitude                           string                 `json:"IpLongitude,omitempty"`                           // Долгота нахождения плательщика
	Issuer                                string                 `json:"Issuer,omitempty"`                                // Название банка-эмитента
	IssuerBankCountry                     string                 `json:"IssuerBankCountry,omitempty"`                     // Двухбуквенный код страны эмитента ISO3166-1
	Description                           string                 `json:"Description,omitempty"`                           // Назначение оплаты
	AuthCode                              string                 `json:"AuthCode,omitempty"`                              // Код авторизации
	Data                                  map[string]interface{} `json:"Data,omitempty"`                                  // Произвольный набор параметров
	Token                                 string                 `json:"Token,omitempty"`                                 // Токен карты для повторных платежей
	TotalFee                              float64                `json:"TotalFee,omitempty"`                              // Значение общей комиссии
	CardProduct                           string                 `json:"CardProduct,omitempty"`                           // Тип карточного продукта
	PaymentMethod                         string                 `json:"PaymentMethod,omitempty"`                         // Метод оплаты (например, T-Pay)
	FallBackScenarioDeclinedTransactionId int64                  `json:"FallBackScenarioDeclinedTransactionId,omitempty"` // Номер первой неуспешной транзакции
	Rrn                                   string                 `json:"Rrn,omitempty"`                                   // Уникальный номер банковской транзакции
	CustomFields                          []interface{}          `json:"CustomFields,omitempty"`                          // Кастомные поля
}
