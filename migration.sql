INSERT INTO users (first_name,last_name,email,mobile,country,state,city,address,password,profile_image,isadmin) VALUES ('Mayur', 'Deshmukh', 'mayur.s.deshmukh1092@gmail.com', '8421985646', 'India', 'M.H', 'Pune', 'ABC', '$2a$08$ue93ZY.v07SRcASU9prhHukjuSiLUShLVB3TnbxSUIkJkX/EFXKAq','assets/users/image.png','t');


INSERT INTO category (cname, description) VALUES ('Clothes','All wearable fabrics, ');
INSERT INTO category (cname, description) VALUES('Electronics',' stores or generates electricity');
INSERT INTO category (cname, description) VALUES('Mobile','The mobile phone can be used to communicate ');
INSERT INTO category (cname, description) VALUES('Watch','A watch is a portable timepiece intended ');
INSERT INTO category (cname, description) VALUES('Books','There are several things to consider in order');
INSERT INTO category (cname, description) VALUES('Sports','Shoes are for regular comfort wear');

-- INSERT INTO category (name, description) VALUES ('Clothes','All wearable fabrics, ');
-- INSERT INTO category (name, description) VALUES('Electronics',' stores or generates electricity');
-- INSERT INTO category (name, description) VALUES('Mobile','The mobile phone can be used to communicate ');
-- INSERT INTO category (name, description) VALUES('Watch','A watch is a portable timepiece intended ');
-- INSERT INTO category (name, description) VALUES('Books','There are several things to consider in order');
-- INSERT INTO category (name, description) VALUES('Sports','Shoes are for regular comfort wear');



-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color, size) VALUES('Polo Shirt', 'Benetton Men  Classic Fit Polo Shirt',511, 10, 5, 10, 1 ,'Polo','Sky Blue', 'Medium');
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color, size) VALUES('Wrangler', 'Men  Slim Fit Jeans', 600, 20, 5, 12, 1 , 'Armani','Charcoal Black','Large');
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color, size) VALUES('Dragon Jacket','Made from the skin of one of the dragons', 700, 40, 5, 9, 1 ,'Veteran','Black','Extra Large');
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color) VALUES('HD Ready Android LED TV ', 'Resolution: HD Ready Android TV (1366x768)', 1200, 20, 10, 12, 2 ,'Samsung','Black');
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color) VALUES('Sony DSC W830 Cyber-Shot 20.1 MP ', 'Shoot Camera (Black) with 8X ', 1500, 50, 15, 8, 2 , 'Samsung','Blue');
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color) VALUES('Charger', 'Mi 10W Charger with Cable (1.2 Meter, Black)', 500, 5, 4, 21, 2 ,'One Plus','White');
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color) VALUES('Apple iPhone 11 Pro (64GB)', '5.8-inch (14.7 cm) ', 60000, 5, 10, 15, 3 ,'Apple', 'Golden');
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color) VALUES('Apple iPhone 11 Pro (64GB) Max', '5.8-inch (14.7 cm) ', 700000, 4, 30, 15, 3, 'Apple','Black' );
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color) VALUES('Apple iPhone XR (64GB)', '6.1-inch (15.5 cm) Liquid Retina HD LCD display',  50000, 6, 20, 15, 3, 'Apple','Grey' );
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color) VALUES ('Rolex Watch','by wearing it you are bound to feel realaxed',2100, 10, 10, 23, 4, 'Rolex', 'Blue');
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color) VALUES ('Titan Watch','With the look of and feel of old days',2000, 15, 8, 5, 4 ,'Titan','Ocean Blue');
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color) VALUES ('Sonata Watch','Stylished belts and longer battery',3010, 20, 12, 5, 4 ,'Sonata','Golden');
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand) VALUES ('Wings of Fire','autobiography by visionary scientist',332, 10, 5, 2, 5, 'TechMax Publications');
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand) VALUES ('Thoughts to Inspire','famous quotes by Swami Vivekananda', 150, 15, 10, 2, 5, 'Technical Publications' );
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand, color, size) VALUES ('Lancer','Mens Running Shoes', 150, 15, 10, 9, 6 , 'Nike', 'Red', 'Small');
-- INSERT INTO products (name, description, price, discount, quantity, tax, category_id, brand) VALUES ('Football','Sporting Goods', 200, 40, 15, 12, 6, 'Cosco' );





-- INSERT INTO productimages VALUES(1,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905019/poloshit1_gnzidq.jpg'); 
-- INSERT INTO productimages VALUES(2,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905299/Wrangler2_l0acph.jpg');
-- INSERT INTO productimages VALUES(3,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905357/Dragon_Jacket1_vh2vkp.jpg');
-- INSERT INTO productimages VALUES(4,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599904937/Android_LED_TV1_ochk0u.jpg');
-- INSERT INTO productimages VALUES(5,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905465/Sony_DSC2_jn3ghu.jpg');
-- INSERT INTO productimages VALUES(6,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905507/Charger1_c0s2a0.jpg');
-- INSERT INTO productimages VALUES(7,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905572/Apple_iPhone_XR2_cynfem.jpg');
-- INSERT INTO productimages VALUES(8,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905572/Apple_iPhone_XR2_cynfem.jpg');
-- INSERT INTO productimages VALUES(9,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905572/Apple_iPhone_XR2_cynfem.jpg');
-- INSERT INTO productimages VALUES(10,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905678/Relax_Watch1_mbgnmj.jpg');
-- INSERT INTO productimages VALUES(11,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905678/Relax_Watch1_mbgnmj.jpg');
-- INSERT INTO productimages VALUES(12,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905678/Relax_Watch1_mbgnmj.jpg');
-- INSERT INTO productimages VALUES(13,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905797/Wings_of_Fire2_ssm4lm.jpg');
-- INSERT INTO productimages VALUES(14,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905822/vivkananadbook1_d2qpiz.jpg');
-- INSERT INTO productimages VALUES(15,'https://res.cloudinary.com/mayur-cloud/image/upload/v1599905887/Mens_Running_Shoes3_krlind.jpg');
-- INSERT INTO productimages VALUES(16, 'https://res.cloudinary.com/mayur-cloud/image/upload/v1600263994/Football3_jab0nt.jpg');





-- INSERT INTO category (cname, description) VALUES ('Clothes','All wearable fabrics, ');
-- INSERT INTO category (cname, description) VALUES('Electronics',' stores or generates electricity');
-- INSERT INTO category (cname, description) VALUES('Mobile','The mobile phone can be used to communicate ');
-- INSERT INTO category (cname, description) VALUES('Watch','A watch is a portable timepiece intended ');
-- INSERT INTO category (cname, description) VALUES('Books','There are several things to consider in order');
-- INSERT INTO category (cname, description) VALUES('Sports','Shoes are for regular comfort wear');



-- INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, size, image_urls) VALUES('Polo Shirt', 'Benetton Men  Classic Fit Polo Shirt',511, 10, 5, 10, 1 ,'Polo','Sky Blue', 'Medium', ARRAY ['https://res.cloudinary.com/mayur-cloud/image/upload/v1599905019/poloshit1_gnzidq.jpg']);
-- INSERT INTO products (name, description, price, discount, quantity, tax, cid, brand, color, size, image_urls) VALUES('Wrangler', 'Men  Slim Fit Jeans', 600, 20, 5, 12, 1 , 'Armani','Charcoal Black','Large', ARRAY ['https://res.cloudinary.com/mayur-cloud/image/upload/v1599905299/Wrangler2_l0acph.jpg']);
