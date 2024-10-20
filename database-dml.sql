INSERT INTO public.categories (user_id,category_type,category_name,category_description,created_at,created_by,updated_at,updated_by) VALUES
	 (2,'IN','Gaji','','2024-10-20 05:16:56.125105','Muhammad Febri Andani',NULL,'system'),
	 (2,'OUT','Makan','','2024-10-20 05:17:49.177416','Muhammad Febri Andani',NULL,'system'),
	 (2,'OUT','Pakaian','','2024-10-20 05:19:44.692571','Muhammad Febri Andani',NULL,'system'),
	 (2,'OUT','Nonton Bioskop','','2024-10-20 05:19:51.935011','Muhammad Febri Andani',NULL,'system'),
	 (2,'IN','Tunjangan','','2024-10-20 05:20:01.020231','Muhammad Febri Andani',NULL,'system'),
	 (2,'IN','Bonus','','2024-10-20 05:20:04.481234','Muhammad Febri Andani',NULL,'system'),
	 (3,'OUT','Sewa Kost Bulanan','','2024-10-20 07:15:04.949879','Muhammad Febri Andani',NULL,'system'),
	 (2,'OUT','Bayar Cicilan','','2024-10-20 08:22:43.596385','Muhammad Febri Andani',NULL,'system');


INSERT INTO public.transactions (user_id,category_id,category_type,amount,description,created_at,created_by,updated_at,updated_by) VALUES
	 (2,1,'IN','7500000','Gaji bulan Oct','2024-10-20 04:59:41.175037','Muhammad Febri Andani',NULL,'system'),
	 (2,2,'OUT','100000','Makan mcd','2024-10-20 05:18:35.869715','Muhammad Febri Andani',NULL,'system'),
	 (2,3,'OUT','2500000','Bayar kost bulanan','2024-10-20 06:37:09.76407','Muhammad Febri Andani',NULL,'system'),
	 (2,10,'OUT','4000000','Bayar Cicilan motor dan lain-lain','2024-10-20 08:23:58.268095','Muhammad Febri Andani',NULL,'system');


INSERT INTO public.users ("name",username,email,phone_number,"password",created_at,created_by,updated_at,updated_by,is_active) VALUES
	 ('Muhammad Febri Andani','febriandani_','febriandani176@gmail.com','62895331276346','$2a$14$MOg75F2PomgA2tj2HmPZK.d02I2k66o6klTptq1LtgotgwjzopGcm','2024-10-20 04:24:56.699311','system','2024-10-20 07:37:42.371995','system',true),
	 ('Muhammad Febri Andani','febriandani18','febriandani@gmail.com','62895331276346','$2a$14$m5vxSNUuydloJgPMrG1SZuFuqcs/ka.4lHncRPZN7czwZ0HGm6Jx2','2024-10-20 03:43:14.851001','system',NULL,'system',true);
