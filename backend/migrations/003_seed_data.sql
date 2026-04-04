-- 003_seed_data.sql

-- Industries (Macedonian names)
INSERT INTO industries (slug, name) VALUES
    ('it_software',        'ИТ и софтвер'),
    ('it_hardware',        'ИТ хардвер и телекомуникации'),
    ('finance_banking',    'Финансии и банкарство'),
    ('manufacturing',      'Производство и индустрија'),
    ('retail_commerce',    'Трговија и продажба'),
    ('healthcare',         'Здравство и фармација'),
    ('education',          'Образование и наука'),
    ('construction',       'Градежништво и недвижнини'),
    ('media_marketing',    'Медиуми, маркетинг и реклама'),
    ('logistics',          'Логистика и транспорт'),
    ('tourism_hospitality','Туризам и угостителство'),
    ('public_sector',      'Јавна администрација'),
    ('ngo_nonprofit',      'НВО и непрофитен сектор'),
    ('legal_consulting',   'Правни и консалтинг услуги'),
    ('agriculture',        'Земјоделство и прехранбена индустрија'),
    ('other',              'Друго')
ON CONFLICT (slug) DO NOTHING;

-- Cities (major Macedonian cities)
INSERT INTO cities (slug, name) VALUES
    ('skopje',       'Скопје'),
    ('bitola',       'Битола'),
    ('kumanovo',     'Куманово'),
    ('prilep',       'Прилеп'),
    ('tetovo',       'Тетово'),
    ('veles',        'Велес'),
    ('stip',         'Штип'),
    ('ohrid',        'Охрид'),
    ('strumica',     'Струмица'),
    ('gostivar',     'Гостивар'),
    ('kavadarci',    'Кавадарци'),
    ('kicevo',       'Кичево'),
    ('struga',       'Струга'),
    ('radovis',      'Радовиш'),
    ('gevgelija',    'Гевгелија'),
    ('kocani',       'Кочани'),
    ('negotino',     'Неготино'),
    ('debar',        'Дебар'),
    ('delcevo',      'Делчево'),
    ('vinica',       'Виница'),
    ('remote',       'Remote')
ON CONFLICT (slug) DO NOTHING;
