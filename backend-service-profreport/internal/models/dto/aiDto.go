package dto

type DiagramStruct struct {
	Name  string `json:"name" description:"Название показателя"`
	Label string `json:"label" description:"Метка показателя"`
	Value int    `json:"value" description:"Значение показателя"`
}

type StructuredResponseSchema struct {
	DiagramRIASEC                         []DiagramStruct `json:"diagram_RIASEC" description:"Таблица показателей RIASEC"`
	DominantRIASEC                        string          `json:"dominant_RIASEC" description:"Дай определение доминирующего кода (например, «S»/«SE»/«SA») и опиши характерные предпочтения, типичные виды деятельности, рабочие среды и риски несоответствия. 4–6 предложений"`
	InterpretationRIASEC                  string          `json:"interpretation_RIASEC" description:"как профиль RIASEC соотносится с образованием и карьерой, где сильные/слабые стороны, какие траектории логичны"`
	DiagramObjectsOfActivityKlimov        []DiagramStruct `json:"diagram_Objects_Of_Activity_Klimov" description:"Таблица показателей объектов деятельности (Климов)"`
	DominantObjectsOfActivityKlimov       string          `json:"dominant_Objects_Of_Activity_Klimov" description:"Определи доминирующий тип объекта труда, опиши подходящие виды деятельности и среды, риски несоответствия. 3–5 предложений."`
	InterpretationObjectsOfActivityKlimov string          `json:"interpretation_Objects_Of_Activity_Klimov" description:"2–4 предложения: как профиль по Климову согласуется с RIASEC и какие типы профессий подходят."`
	DiagramPersonalQualities              []DiagramStruct `json:"diagram_Personal_Qualities" description:"Таблица показателей личностных качеств"`
	DominantPersonalQualities             string          `json:"dominant_Personal_Qualities" description:"Выдели и опиши наиболее выраженные качества (порог — верхний квантиль или балл ≥ 4 из 5), приведи 6–10 пунктов как markdown список"`
	InterpretationPersonalQualities       string          `json:"interpretation_Personal_Qualities" description:"2–4 предложения: как качества влияют на выбор профессий/сред, где сильные стороны полезны, где возможны риски."`
	IntegratedResultsAnalysis             string          `json:"integrated_results_analysis" description:"Сопоставление результатов по всем четырём блокам (RIASEC, Климов, качества, ценности): согласованность профилей; качества, поддерживающие успех в релевантных профессиях; ключевые ценности и их связь с рекомендуемыми направлениями; итог по склонностям/сильным сторонам/ценностным ориентациям. 3–5 абзацев"`
}
