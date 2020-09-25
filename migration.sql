INSERT INTO category (cname, description) VALUES ('Clothes','All wearable fabrics, ');
INSERT INTO category (cname, description) VALUES('Electronics',' stores or generates electricity');
INSERT INTO category (cname, description) VALUES('Mobile','The mobile phone can be used to communicate ');
INSERT INTO category (cname, description) VALUES('Watch','A watch is a portable timepiece intended ');
INSERT INTO category (cname, description) VALUES('Books','There are several things to consider in order');
INSERT INTO category (cname, description) VALUES('Sports','Shoes are for regular comfort wear');



INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, size, image_urls) VALUES('Polo Shirt', 'Benetton Men  Classic Fit Polo Shirt',511, 10, 5, 10, 1 ,'Polo','Sky Blue', 'Medium', ARRAY ['poloshit1.jpeg','poloshirt2.jpeg', 'poloshirt3.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, size, image_urls) VALUES('Wrangler', 'Men  Slim Fit Jeans', 600, 20, 5, 12, 1 , 'Armani','Charcoal Black','Large', ARRAY ['Wrangler2','Wrangler3']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, size, image_urls) VALUES('Dragon Jacket','Made from the skin of one of the dragons', 700, 40, 5, 9, 1 ,'Veteran','Black','Extra Large', ARRAY ['DragonJacket1.jpeg','DragonJacket2.jpeg', 'DragonJacket3.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, image_urls) VALUES('HD Ready Android LED TV ', 'Resolution: HD Ready Android TV (1366x768)', 1200, 20, 10, 12, 2 ,'Samsung','Black', ARRAY ['AndroidLEDTV1.jpeg', 'AndroidLEDTV2.jpeg','AndroidLEDTV3.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, image_urls) VALUES('Sony DSC W830 Cyber-Shot 20.1 MP ', 'Shoot Camera (Black) with 8X ', 1500, 50, 15, 8, 2 , 'Samsung','Blue', ARRAY ['SonyDSC1.jpeg','SonyDSC2.jpeg', 'SonyDSC3.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, image_urls) VALUES('Charger', 'Mi 10W Charger with Cable (1.2 Meter, Black)', 500, 5, 4, 21, 2 ,'One Plus','White', ARRAY ['Charger1.jpeg','Charger2.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, image_urls) VALUES('Apple iPhone 11 Pro (64GB)', '5.8-inch (14.7 cm) ', 60000, 5, 10, 15, 3 ,'Apple', 'Golden', ARRAY ['AndroidLEDTV1.jpeg', 'AndroidLEDTV2.jpeg','AndroidLEDTV3.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, image_urls) VALUES('Apple iPhone 11 Pro (64GB) Max', '5.8-inch (14.7 cm) ', 700000, 4, 30, 15, 3, 'Apple','Black' , ARRAY ['AndroidLEDTV1.jpeg', 'AndroidLEDTV2.jpeg','AndroidLEDTV3.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, image_urls) VALUES('Apple iPhone XR (64GB)', '6.1-inch (15.5 cm) Liquid Retina HD LCD display',  50000, 6, 20, 15, 3, 'Apple','Grey' , ARRAY ['AndroidLEDTV1.jpeg', 'AndroidLEDTV2.jpeg','AndroidLEDTV3.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, image_urls) VALUES ('Rolex Watch','by wearing it you are bound to feel realaxed',2100, 10, 10, 23, 4, 'Rolex', 'Blue', ARRAY ['RelaxWatch1.jpeg','RelaxWatch2.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, image_urls) VALUES ('Titan Watch','With the look of and feel of old days',2000, 15, 8, 5, 4 ,'Titan','Ocean Blue', ARRAY ['TitanWatch1.jpeg','TitanWatch2.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, image_urls) VALUES ('Sonata Watch','Stylished belts and longer battery',3010, 20, 12, 5, 4 ,'Sonata','Golden', ARRAY ['SonataWatch1.jpeg','SonataWatch2.jpeg', 'SonataWatch3.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, image_urls) VALUES ('Wings of Fire','autobiography by visionary scientist',332, 10, 5, 2, 5, 'TechMax Publications', ARRAY ['WingsofFire1.jpeg','WingsofFire2.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, image_urls) VALUES ('Thoughts to Inspire','famous quotes by Swami Vivekananda', 150, 15, 10, 2, 5, 'Technical Publications' , ARRAY ['vivkananadbook1.jpeg','vivkananadbook2.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, size, image_urls) VALUES ('Lancer','Mens Running Shoes', 150, 15, 10, 9, 6 , 'Nike', 'Red', 'Small', ARRAY ['MensRunningShoes1.jpeg','MensRunningShoes2.jpeg']);
INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand) VALUES ('Football','Sporting Goods', 200, 40, 15, 12, 6, 'Cosco' );
