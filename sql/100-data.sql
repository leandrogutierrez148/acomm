-- Insert 3 product categories
INSERT INTO product_categories (code, name)
VALUES ('CLTH', 'Clothing'),
    ('SHS', 'Shoes'),
    ('ACC', 'Accessories');
-- Insert Brands
INSERT INTO brands (
        name,
        meta_tag_title,
        meta_tag_description,
        sort_order
    )
VALUES (
        'Nike',
        'Nike Shoes and Apparel',
        'Buy Nike products online',
        1
    ),
    (
        'Adidas',
        'Adidas Gear',
        'Premium Adidas clothing',
        2
    ),
    ('Puma', 'Puma Store', 'Puma athletic wear', 3);
-- Insert 3 VTEX products (ids 1, 2, 3)
INSERT INTO products (
        name,
        department_id,
        category_id,
        brand_id,
        link_id,
        ref_id,
        description,
        keywords,
        title,
        is_active,
        score
    )
VALUES (
        'Zapatillas adidas Runfalcon 5 De Hombre',
        106,
        2,
        2,
        'zapatillas-adidas-runfalcon-5-de-hombre--6ih7758-000',
        '6ID8760',
        '<html><head><meta charset="UTF-8"></head><body><span><b>Zapatillas de running para uso diario que te ayudan a llegar hasta el final</b><br><br>  Tanto en la pista como en la cinta de correr, alcanzá todas tus metas con estas zapatillas de running adidas. Incorporan una mediasuela con amortiguación Cloudfoam que te ofrece una pisada más cómoda y suave. La parte superior de malla transpirable y la suela Adiwear de gran resistencia al desgaste la convierten en una silueta perfecta para llevar durante todo el día.<br><br><b>DETALLES</b><br><br><li> Horma clásica.<br><li> Sistema de atado de cordones.<br><li> Parte superior de malla.<br><li> Forro interno textil.<br><li> Plantilla OrthoLite®.<br><li> Mediasuela Cloudfoam.<br><li> Caída mediasuela: 10 mm (talón: 33 mm / antepié: 23 mm).<br><li> Suela Adiwear.<br><li> Color del artículo: Core Black / Cloud White / Core Black.</span></body></html>',
        '6ID8760,running',
        'Zapatillas adidas Runfalcon 5 De Hombre',
        true,
        150
    ),
    (
        'Camiseta adidas Copa Mundial Argentina 26 Local Mujer',
        106,
        1,
        2,
        'camisetas-de-equipos-adidas-camiseta-titular-seleccion-argentina-26-mujer-59140',
        '6IH8218',
        'Celebrá el legado futbolístico de la Albiceleste con la camiseta titular de la Selección Argentina 26. La camiseta rinde homenaje a los Campeones del Mundo, mostrando un efecto degradado inspirado en las camisetas usadas durante las tres victorias de Argentina en la Copa del Mundo.',
        'camiseta,argentina',
        '',
        true,
        0
    ),
    (
        'Gorra Puma Essentials N°1 Logo Patch Unisex',
        106,
        3,
        3,
        'gorra-puma-essentials-n1-logo-patch-unisex-02599711-000',
        '02599711',
        '<html><head><meta charset="UTF-8"></head><body><span>Lucí un look clásico en esta gorra desestructurada, que presenta el logo PUMA. La visera curvada te protege del sol y la banda interna absorbe el sudor, para que puedas mantener tu concentración.</span></body></html>',
        '02599711,moda',
        'Gorra Puma Essentials N°1 Logo Patch Unisex',
        true,
        80
    );
-- Zapatillas (1 Color, 3 Talles) -> Items 1, 2, 3
INSERT INTO items (
        product_id,
        sku,
        price,
        is_active,
        activate_if_possible,
        name,
        ref_id,
        packaged_height,
        packaged_length,
        packaged_width,
        packaged_weight_kg,
        cubic_weight,
        is_kit,
        manufacturer_code,
        commercial_condition_id,
        measurement_unit,
        unit_multiplier,
        kit_itens_sell_apart
    )
VALUES (
        1,
        '6ID8760-000-39',
        120.99,
        true,
        true,
        '39',
        '6ID8760-000-39',
        10,
        10,
        10,
        1,
        0.2083,
        false,
        '6ID8760',
        9,
        'un',
        1,
        false
    ),
    (
        1,
        '6ID8760-000-40',
        120.99,
        true,
        true,
        '40',
        '6ID8760-000-40',
        10,
        10,
        10,
        1,
        0.2083,
        false,
        '6ID8760',
        9,
        'un',
        1,
        false
    ),
    (
        1,
        '6ID8760-000-41',
        120.99,
        true,
        true,
        '41',
        '6ID8760-000-41',
        10,
        10,
        10,
        1,
        0.2083,
        false,
        '6ID8760',
        9,
        'un',
        1,
        false
    );
-- Camiseta (3 Talles) -> Items 4, 5, 6
INSERT INTO items (
        product_id,
        sku,
        price,
        is_active,
        activate_if_possible,
        name,
        ref_id,
        packaged_height,
        packaged_length,
        packaged_width,
        packaged_weight_kg,
        cubic_weight,
        is_kit,
        manufacturer_code,
        commercial_condition_id,
        measurement_unit,
        unit_multiplier,
        kit_itens_sell_apart
    )
VALUES (
        2,
        '6IH8218-S',
        89.99,
        true,
        true,
        'S',
        '6IH8218-S',
        10,
        10,
        10,
        0.5,
        0.1,
        false,
        '6IH8218',
        9,
        'un',
        1,
        false
    ),
    (
        2,
        '6IH8218-M',
        89.99,
        true,
        true,
        'M',
        '6IH8218-M',
        10,
        10,
        10,
        0.5,
        0.1,
        false,
        '6IH8218',
        9,
        'un',
        1,
        false
    ),
    (
        2,
        '6IH8218-L',
        89.99,
        true,
        true,
        'L',
        '6IH8218-L',
        10,
        10,
        10,
        0.5,
        0.1,
        false,
        '6IH8218',
        9,
        'un',
        1,
        false
    );
-- Gorra (2 Colores, 1 Talle) -> Items 7, 8
INSERT INTO items (
        product_id,
        sku,
        price,
        is_active,
        activate_if_possible,
        name,
        ref_id,
        packaged_height,
        packaged_length,
        packaged_width,
        packaged_weight_kg,
        cubic_weight,
        is_kit,
        manufacturer_code,
        commercial_condition_id,
        measurement_unit,
        unit_multiplier,
        kit_itens_sell_apart
    )
VALUES (
        3,
        '02599711-000-U',
        25.00,
        true,
        true,
        'U',
        '02599711-000-U',
        5,
        5,
        5,
        0.2,
        0.05,
        false,
        '02599711',
        9,
        'un',
        1,
        false
    ),
    (
        3,
        '02599711-001-U',
        25.00,
        true,
        true,
        'U',
        '02599711-001-U',
        5,
        5,
        5,
        0.2,
        0.05,
        false,
        '02599711',
        9,
        'un',
        1,
        false
    );
-- Insert sample specifications for product 1
INSERT INTO specifications (product_id, key, value)
VALUES (1, 'Material', 'Sintetico'),
    (1, 'Color', 'Negro');
-- ITEM IMAGES
-- item = 1 indicates SKU 6ID8760-000-39 (We attach the exact VTEX images provided)
INSERT INTO item_images (
        item_id,
        archive_id,
        sku_id,
        name,
        is_main,
        text,
        label,
        url,
        file_location,
        position
    )
VALUES (
        1,
        1593387,
        282421,
        'Shein 1',
        true,
        'Shein 1',
        'Detalle 1',
        'https://img.ltwebstatic.com/images3_spmp/2025/01/14/07/1736785958fe618e066f77a9d15a14f91d9bb62bba_thumbnail_900x.webp',
        '',
        0
    ),
    (
        1,
        1593396,
        282421,
        'Shein 2',
        false,
        'Shein 2',
        'Detalle 2',
        'https://img.ltwebstatic.com/images3_spmp/2025/01/14/64/1736785962f5139f2a349ebf67b526159b6f257938_thumbnail_900x.webp',
        '',
        1
    ),
    (
        1,
        1593409,
        282421,
        'Shein 3',
        false,
        'Shein 3',
        'Detalle 3',
        'https://img.ltwebstatic.com/images3_spmp/2025/01/14/7e/17367859646539ad10000cc8e6a2402a87fa490e13_thumbnail_900x.webp',
        '',
        2
    ),
    (
        1,
        1593417,
        282421,
        'Shein 4',
        false,
        'Shein 4',
        'Detalle 4',
        'https://img.ltwebstatic.com/images3_spmp/2025/01/14/6f/17367859603c73085a4e889346b2234a0343becd8e_thumbnail_900x.webp',
        '',
        3
    );
-- item = 4 indicates first Camiseta SKU (6IH8218-S).
INSERT INTO item_images (
        item_id,
        archive_id,
        sku_id,
        name,
        is_main,
        text,
        label,
        url,
        file_location,
        position
    )
VALUES (
        4,
        2263274,
        323003,
        'Adidas 1',
        true,
        'Adidas 1',
        'Detalle 1',
        'https://assets.adidas.com/images/h_2000,f_auto,q_auto,fl_lossy,c_fill,g_auto/841749a208934f2ab4a0cfa1a8ae237d_9366/Camiseta_primera_equipacion_Argentina_26_Blanco_JM8396_21_model.jpg',
        '',
        0
    ),
    (
        4,
        2263275,
        323003,
        'Adidas 2',
        false,
        'Adidas 2',
        'Detalle 2',
        'https://assets.adidas.com/images/h_2000,f_auto,q_auto,fl_lossy,c_fill,g_auto/3c16ce3f59834e0cbe052a59f63e27c5_9366/Camiseta_primera_equipacion_Argentina_26_Blanco_JM8396_23_hover_model.jpg',
        '',
        1
    ),
    (
        4,
        2263276,
        323003,
        'Adidas 3',
        false,
        'Adidas 3',
        'Detalle 3',
        'https://assets.adidas.com/images/h_2000,f_auto,q_auto,fl_lossy,c_fill,g_auto/6e10e88e911b47a6ad914e3626ea6b67_9366/Camiseta_primera_equipacion_Argentina_26_Blanco_JM8396_25_model.jpg',
        '',
        2
    ),
    (
        4,
        2263277,
        323003,
        'Adidas 4',
        false,
        'Adidas 4',
        'Detalle 4',
        'https://assets.adidas.com/images/h_2000,f_auto,q_auto,fl_lossy,c_fill,g_auto/20f23eb7f4df4dd4958601c064ffa0d8_9366/Camiseta_primera_equipacion_Argentina_26_Blanco_JM8396_01_laydown.jpg',
        '',
        3
    );
-- item = 7 indicates Gorra (02599711-000-U).
INSERT INTO item_images (
        item_id,
        archive_id,
        sku_id,
        name,
        is_main,
        text,
        label,
        url,
        file_location,
        position
    )
VALUES (
        7,
        2580631,
        310272307,
        'Puma 1',
        true,
        'Puma 1',
        'Detalle 1',
        'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/025997/01/fnd/EEA/fmt/png/Gorra-de-b%C3%A9isbol-ESS-Patch-con-logotipo-n.%C2%BA-1',
        '',
        0
    ),
    (
        7,
        2580632,
        310272307,
        'Puma 2',
        false,
        'Puma 2',
        'Detalle 2',
        'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/025997/01/mod01/fnd/EEA/fmt/png/Gorra-de-b%C3%A9isbol-ESS-Patch-con-logotipo-n.%C2%BA-1',
        '',
        1
    ),
    (
        7,
        2580633,
        310272307,
        'Puma 3',
        false,
        'Puma 3',
        'Detalle 3',
        'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/025997/01/mod02/fnd/EEA/fmt/png/Gorra-de-b%C3%A9isbol-ESS-Patch-con-logotipo-n.%C2%BA-1',
        '',
        2
    ),
    (
        7,
        2580634,
        310272307,
        'Puma 4',
        false,
        'Puma 4',
        'Detalle 4',
        'https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/025997/01/bv/fnd/EEA/fmt/png/Gorra-de-b%C3%A9isbol-ESS-Patch-con-logotipo-n.%C2%BA-1',
        '',
        3
    );