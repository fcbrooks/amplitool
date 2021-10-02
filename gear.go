package main

var Amps = map[string]string{
	// Clean
	"5f4f50a1-d5cb-43be-ad11-084e4ff21ea6": "'57 Champ",
	"4c9e667b-932a-42e3-a5d8-a9d9374c9959": "'57 Custom Twin-Amp",
	"f0951b1e-91d2-4360-80d7-793fa785d2d6": "'64 Vibroverb Custom",
	"89b3caab-dffb-4c29-85d9-2a60cb93c566": "'65 Deluxe Reverb",
	"2a1f483c-a136-45b6-81ec-e92c60f8d009": "'65 Princeton",
	"d3c791b9-58f1-41d2-8a88-797e98cc5b29": "'65 Super Reverb",
	"b3869f27-a9f1-4482-add4-9512c16917ea": "'65 Twin Reverb",
	"a91067a3-fd80-40a8-be35-0681da5c4f47": "American Clean MKIII",
	"71a76a9f-cf70-4f59-971f-9864a055523c": "American Tube Clean 1",
	"82972243-cd55-4b43-82f3-f15e3bc13dc7": "American Tube Clean 2",
	"ca4587b9-3960-49de-9509-5a61e9b5cbae": "American Vintage B",
	"84f03443-ae64-4c7e-970f-06d1191cd906": "American Vintage D",
	"95b9bc84-89fa-48f5-a336-26a30a044ca3": "American Vintage T",
	"016a8c2a-489e-49da-81d7-5b72feb60f74": "Champion 600",
	"48503c68-f5e4-40d6-a4e1-75b0168e5e6f": "Custom Solid State Clean",
	"ac08939a-32bf-496c-96ac-5d6c530abf14": "Jazz Amp 120",
	"300ed819-d21b-4589-b095-afd038e9f08c": "JH 1200",
	"abdcae70-bff2-4b02-bf2f-d716dd8e8adf": "MAZ 18 Jr",
	"15761216-f2fe-4d41-a6ec-9bff8199517c": "Metal Clean T",
	"d0546d04-505c-42b1-8e9e-668a16adcfa8": "Pro Junior",
	"a2f18e96-4d56-4372-b438-11bd0f42f6f3": "SilverTwelve",
	"dffa559d-7b12-464a-9fbf-877ca25f5cf3": "Vibro-King",
	// Crunch
	"f2190a68-52ea-408a-9c39-2ea8279c0d43": "'53 Bassman",
	"6f1c22b5-3593-4d86-a9a3-fae8c9504d77": "'57 Bandmaster",
	"a0fa7c56-0772-4ddd-9320-c2ee254a3c4a": "'57 Custom Champ",
	"bf860ad9-cd8a-425b-8049-29211fce237a": "'57 Custom Deluxe",
	"6c421302-9602-4ee8-b94a-672aa24cdde4": "'57 Custom Pro-Amp",
	"d4d5b530-0ce1-46cf-a47e-bf0224fa715e": "'57 Deluxe",
	"3fcc8ad1-6d5e-416d-9c3d-7aae91c6f4d4": "'59 Bassman LTD",
	"dd7b0e06-a17a-4851-83c4-ee32ca303b01": "AD 30",
	"0d4c8b80-92d6-4f40-8178-51ba0179eb1d": "American Tube Vintage",
	"f058124b-498f-4899-8b29-35453d6aecff": "Bi-Valve",
	"533d3c6c-b3cd-455c-a3a1-642016f5cda9": "British Blue Tube 30TB",
	"5d235e0d-9fd7-429e-b483-6f815281f3d7": "British Copper 30TB",
	"d089ef66-b5c4-4274-910c-6a6ee194cf04": "British Lead S100",
	"2a95b351-ba28-473d-b0a7-fd924f32d9f9": "Custom Solid State Fuzz",
	"3c25674f-a418-4fec-863c-f94495c746a0": "Dual Terror",
	"26fbbf20-f88e-46de-a76f-5aabd2c8fd8d": "HiAmp",
	"7788f707-4ef2-44cd-862a-a82ffdf7172b": "JH Gold",
	"827aedfb-cdc1-412e-8e47-5bac3c3c6d06": "OR 50",
	"6e8690b3-f6cf-4c36-b3c2-7f38fcc5706e": "OR-120",
	"e1eed2cf-6777-46c4-ada2-65df0d7afc46": "Red Pig",
	"f4b89ab3-8ca6-44ee-b90b-a570040c8a3d": "Super-Sonic",
	"99e446c7-49df-45b1-bff9-26d95e10c763": "Tiny Terror",
	"1284c9cc-6efa-4720-a0da-106a2d2af1d8": "TransAtlantic TA-30",
	"bd903ed7-82bf-41cc-9834-685b6e3667b5": "Tube Vintage Combo",
	"cc59b472-2a2f-40b1-97f4-6ee4b7536c87": "Z Wreck",
	// High Gain
	"2ea3ecfb-1b0c-417a-8788-86f5915f43c5": "AFD 100",
	"4af9d89a-c06b-4c8f-b137-af72bc58fded": "American Lead MkIII",
	"8fe96936-5178-4950-9b80-d89c32534bad": "Brit 8000",
	"cbf3c00f-dc31-4c7f-a409-f7fdbca005a8": "Brit 9000",
	"3930eb8b-3eda-4079-b86d-7bfd7d4449bc": "Brit Silver",
	"f970b981-527b-4eb8-ab92-fee301d74678": "Brit Valve Pre",
	"fb5fc82f-a926-4591-87d2-168906fd79d3": "British Tube Lead 1",
	"12db8dc3-fbda-478d-98f8-64ce892478d5": "British Tube Lead 2",
	"cd657b5a-7cc1-4934-b296-58789188b662": "Custom Modern Hi-Gain",
	"2ed045ad-a344-4b35-b95e-0d3a3c1220ce": "Custom Solid State Lead",
	"75ad4a0e-5c75-443d-8617-9681c4fe58d3": "Dual Rectifier",
	"c936fc9c-1594-48c8-b561-824827452a66": "E650",
	"55078333-dcfd-41ac-9e87-cd6ea334507a": "JCA20H",
	"913db945-c3a9-4e96-ad17-9c8d1053d913": "JCA100H",
	"155f0121-a2ee-4e16-aaa0-44948f9be44f": "JCM Slash",
	"6ec4bf7a-dc59-4443-b2fb-1e645bf5192c": "Mark III",
	"1fbf7d6e-dad8-470f-b204-4d96b5466893": "Mark IV",
	"9400d18f-5f72-40ac-aa37-861ba3f18da5": "Metal Lead T",
	"dcc7c825-76f4-4703-8e1f-b8a12b30b1de": "Metal Lead V",
	"802af8da-63c6-4ccf-a4b9-6d13255ef57f": "Metal Lead W",
	"907be0ce-a419-4281-901f-dcd6763de54a": "MH-500 Metalhead",
	"bca11751-a7c1-49f5-846d-031f7eb780f0": "Modern Tube Lead",
	"88d927a0-e399-4a1d-ac68-0699eee85f02": "Powerball",
	"e6151532-1028-422c-9a5d-fc57594ce8e8": "RockerVerb 50",
	"4a22ac9f-aabb-4180-b697-5d5710a1acc2": "SLO 100",
	"e3260631-d81f-4c76-9e4f-d12be6ede5cb": "Thunderverb 200",
	"c85e5dc4-d051-4aad-846f-038b0b5233c5": "Triple Rectifier",
	"5558e374-6d37-4674-a05e-2e005830d24e": "V3M",
	"1b5961b1-f862-4c8a-9a9b-a920da8c5cc2": "Vintage Metal Lead",
	// Bass
	"ecb60014-617c-4637-9435-28c1480a0e8f": "360Bass Preamp",
	"dfb00647-6603-4fe1-a67a-5690a4dad0fb": "AD 200",
	"9e6f407a-161d-433b-bddc-8565103fc9ce": "Bassman 300",
	"9d71083c-4f67-4e3d-9a1a-77431a6d1c10": "Green BA250",
	"d33c157e-aa68-4cba-a9cb-4e2cff5c3caf": "MB 150 S",
	"ad4ea282-ced9-49d0-9670-e9782ce5c5b7": "Solid State Bass Preamp",
	"0265b273-d648-47c7-a5ef-579acba82a0a": "SVX-4B",
	"41f3868c-62c3-4bd0-8c29-130e9426d4e9": "SVX-15N",
	"ff274db4-43c3-4fb9-b44d-d04aefd13b28": "SVX-15R",
	"f1d8d4c0-770c-469e-88fc-0fa2ffe7e8bc": "SVX-500",
	"52f28b23-80e3-4f43-9508-4447258b11c0": "SVX-CL",
	"862b5977-9c12-4665-88cd-86f668da8877": "SVX-Pro",
	"2aa0f50f-a6c9-4edd-97c2-df71a24087db": "SVX-VR",
	"131366d6-a73d-43bb-8357-d1f7b78b79a5": "TBP-1",
}

var Cabs = map[string]string{
	// 1x6
	"4b4c561b-68d7-4311-ae31-2432817850bd": "Champion 600",
	// 1x8
	"36ca806c-5136-4370-baa4-ce45b3c1c9af": "1x8 '57 Custom Champ",
	"f7f9e974-28f8-4303-b177-39e620c60e9e": "'57 Champ",
	// 1x10
	"223ef7e4-1afe-4c89-b707-362f0c100d04": "'65 Princeton",
	"06834dae-1774-4a8a-9ff0-e8a1f05b4ae5": "Pro Junior",
	// 1x12
	"7efbb008-6942-4538-875b-d65fc7321617": "1x12 '57 Custom Deluxe",
	"9644a358-408c-4b1e-9948-806e19e076f3": "1x12 Combo",
	"82409e48-b67c-4762-8256-963f43240ccc": "1x12 Mark III",
	"8c6c0893-9f50-492c-ac1d-4773fd697638": "1x12 Mark IV",
	"c895007f-963a-499c-bd76-7cb9fbf31a96": "1x12 MAZ 18 Jr",
	"6252f4e8-9ab0-4283-aa5f-aca9f09309ff": "1x12 MB 150 S",
	"736a416f-dc96-4781-9ef3-5559fd2e7d17": "1x12 MB II",
	"46e8ee66-4ff8-44c6-945d-3a2dfe1323e9": "1x12 MB III",
	"7b8ece5b-4912-4b3d-ad88-05e43245d7fc": "1x12 Open Modern",
	"bcec521d-c918-4af6-98be-b50b644ac3dd": "1x12 Open Vintage",
	"e49d0a85-bcde-4d98-801c-24166a47c8a1": "1x12 PPC 112",
	"33238dfe-a58a-4b2c-a5a7-23388006e414": "1x12 Tiny Terror",
	"70b17311-7b73-4415-968e-c26543e22e18": "'57 Deluxe",
	"57aa5488-221b-4f3f-bb80-bdf0083efe52": "'65 Deluxe Reverb",
	"67012b6a-0886-41d3-9a03-daa13ed47c34": "Super-Sonic",
	// 1x15
	"55ca6af1-4c70-44e8-957e-4f07cc6861b9": "1x15  Bass Vintage",
	"0c2f1129-e09d-4510-9624-7c8b05f188cf": "1x15  OBC 115",
	"151ace02-771b-44d1-afea-0ca3056c3b2c": "1x15 '53 Bassman",
	"59f2b03d-40c1-405b-aeea-38df617b49fa": "1x15 '57 Custom Pro-Amp",
	"7a1b419a-fbd7-4507-93e3-d6dba20ea7ad": "'64 Vibroverb Custom",
	"4da81432-6ff0-4ed4-9df0-119b19ea8681": "SVX-15N",
	"291a5adb-cc71-4eda-be0d-5a0af6bea973": "SVX-15R",
	// 1x18
	"8f3296c5-9056-4fe5-b9c6-d0d0db113080": "1x18 Horn Bass",
	// 2x10
	"63387c79-e174-4d78-9d5e-86fdaa8bc45e": "SVX-500",
	// 2x12
	"2af64b91-de12-48cf-a90e-2e5c4bcc9530": "2x12 '57 Custom Twin-Amp",
	"1d1c171f-7380-4727-8bd6-c63021db9dd6": "2x12 AD 30",
	"8e98245b-fd88-48b4-8951-2d6171c717ce": "2x12 Closed Vintage",
	"3c9256bb-9b75-4dd2-8d49-c9236eedfcb9": "2x12 Gry British Vint",
	"f7902634-12e9-4a2d-9f9a-bcd22781cdab": "2x12 JP Jazz",
	"3e1e42ab-294b-4816-953e-ef800bfdbbe6": "2x12 Open SL",
	"afe74c01-d70d-47ba-bfd5-4ac43ecaa954": "2x12 Open T J120",
	"d9083a81-5a8b-4b19-990e-4d933bed5067": "2x12 Open Vintage",
	"16b2e136-c230-4598-9931-5b913efb76e6": "2x12 PPC 212",
	"f5e43052-b042-4f7c-b49d-a96c6571a4ce": "2x12 PPC OB",
	"d004a94e-14b6-4e92-80eb-89aa306f0008": "2x12 Recifier Horizontal",
	"356da432-3493-4403-ac2d-761a28855638": "2x12 TransAtlantic TA-30",
	"b1de33d6-c660-4bf6-b290-79626a35149e": "2x12 V3M",
	"bce514e5-5fa0-40c3-ae1d-85093a3d62f6": "2x12 Z Wreck",
	"8b37839c-5798-4584-8aca-b7bfce4819a6": "'65 Twin Reverb",
	"a56735e4-226b-4887-bbb7-9a04fa080c0f": "SVX-212 AV",
	"6ac86da3-7341-4a2f-b9c4-ecd118c559c4": "SVX-212H",
	// 2x15
	"17d311d6-440e-4d5b-a695-f74d55e986dd": "2x15 Closed B J130",
	"a55c9112-55bf-4fc2-845d-fd4ce59eb40c": "2x15 Closed D J130",
	// 3x10
	"21a4eb35-fce4-482a-9199-ace3f4205ede": "3x10 '57 Bandmaster",
	"104346d1-5fbc-474d-bcf0-7980cbe04a82": "Vibro-King",
	// 4x10
	"fa7d0a6a-23e9-4f05-a062-607f5cfe7dda": "4x10 '65 Super Reverb",
	"c4af1d27-afd0-426f-9c96-e2f65378b93a": "4x10 OBC 410",
	"9fa8c924-6543-4085-b55b-58b99aada17e": "4x10 Open Vintage",
	"6ea64fbe-8fd5-4eac-a216-28cb0f07faa3": "4x10+tw Bass",
	"3b999e93-6120-4a91-b3d7-0f2902be747f": "4x10+tw TE Bass",
	"4614704e-7ca2-4736-a750-648fe9033650": "'59 Bassman",
	"fad46e0e-3760-484f-9f7b-85f0376780ef": "SVX-410B",
	"1052d32a-7bd7-4463-98c0-b9cc18a000be": "SVX-410S",
	// 4x12
	"936efc52-2172-4faf-9be5-e8f45244a2b9": "4x12 196BV SL",
	"c97bc69c-c02d-4cce-b19d-859b72833550": "4x12 1960AV SL",
	"cd231131-d053-4193-bd88-7bbd14144680": "4x12 Brit 30",
	"7c0b8ce1-cbb4-4e5b-9973-a572143ddb2b": "4x12 Brit 8000",
	"a54d97da-cd7f-4742-acba-45cb2688c8c9": "4x12 Brit 9000",
	"c6dc5147-0436-482f-9a8d-070ccea23c46": "4x12 Brit Silver",
	"79818fe5-e89f-437b-86bc-de80ecea216d": "4x12 Closed 25 C",
	"c4ea21cc-6444-4779-9eee-62d4bc085410": "4x12 Closed 75 C",
	"9ef10de4-2781-4ab7-9179-c1cc4b6615e5": "4x12 Closed HiAmp",
	"0f086cfd-b793-4b69-894a-b69d4c32154b": "4x12 Closed J120",
	"67f95a0d-34e8-4206-b321-3e57c8d1b407": "4x12 Closed Modern",
	"445c0a64-c729-4502-b7c4-91e211a7fc21": "4x12 Closed Vintage",
	"8a9bd7a7-d080-4023-9a60-e8cfcddccf0e": "4x12 Metal F 1",
	"8b147712-d44d-4564-a40f-fe288110ea6c": "4x12 Metal T 1",
	"2bca68bf-1cbd-4b37-b6c2-1b40e1902e5c": "4x12 Metal V 1",
	"ff9cad13-059d-46bf-815a-d8f851c410d4": "4x12 Modern M 1",
	"4edd00f5-1dc0-4130-a3bb-7eb608e834ab": "4x12 PPC 412",
	"849b3340-9e28-411f-9faf-e99b7b2bfb36": "4x12 Recto Traditional Slant",
	"6dfb576f-b549-4dd4-ac79-1021bcb53bb2": "4x12 Red Pig",
	"81866d8a-1072-4bcd-adce-861f842f4e5c": "4x12 Vintage M 1",
	"8eaae0e6-4c9b-471e-814b-6a4596a4d153": "E 412 PRO XXL",
	"088fe7ff-34f6-4a64-9a7c-0f938f05553d": "E 412 Standard",
	// 8x10
	"fb54d283-dfc2-402e-8c76-87713387d770": "8x10  OBC 810",
	"b3900806-4600-4cc0-a2f2-217b9ec09c5f": "Bass 810 Pro",
	"b31f4357-3c42-4c89-8736-64f25bcbed9d": "SVX-810 AV",
	"b11418d6-9a01-42a9-a85b-cd287ac7a98e": "SVX-810E",
	// Other
	"846d771b-3037-4d47-90c7-30e33e32af6f": "122",
	"c275ee33-8180-4b0c-9215-dd6f37aa2394": "Vibratone",
}

var SpeakerCount = map[string]int{
	// 1x6
	"4b4c561b-68d7-4311-ae31-2432817850bd": 1,
	// 1x8
	"36ca806c-5136-4370-baa4-ce45b3c1c9af": 1,
	"f7f9e974-28f8-4303-b177-39e620c60e9e": 1,
	// 1x10
	"223ef7e4-1afe-4c89-b707-362f0c100d04": 1,
	"06834dae-1774-4a8a-9ff0-e8a1f05b4ae5": 1,
	// 1x12
	"7efbb008-6942-4538-875b-d65fc7321617": 1,
	"9644a358-408c-4b1e-9948-806e19e076f3": 1,
	"82409e48-b67c-4762-8256-963f43240ccc": 1,
	"8c6c0893-9f50-492c-ac1d-4773fd697638": 1,
	"c895007f-963a-499c-bd76-7cb9fbf31a96": 1,
	"6252f4e8-9ab0-4283-aa5f-aca9f09309ff": 1,
	"736a416f-dc96-4781-9ef3-5559fd2e7d17": 1,
	"46e8ee66-4ff8-44c6-945d-3a2dfe1323e9": 1,
	"7b8ece5b-4912-4b3d-ad88-05e43245d7fc": 1,
	"bcec521d-c918-4af6-98be-b50b644ac3dd": 1,
	"e49d0a85-bcde-4d98-801c-24166a47c8a1": 1,
	"33238dfe-a58a-4b2c-a5a7-23388006e414": 1,
	"70b17311-7b73-4415-968e-c26543e22e18": 1,
	"57aa5488-221b-4f3f-bb80-bdf0083efe52": 1,
	"67012b6a-0886-41d3-9a03-daa13ed47c34": 1,
	// 1x15
	"55ca6af1-4c70-44e8-957e-4f07cc6861b9": 1,
	"0c2f1129-e09d-4510-9624-7c8b05f188cf": 1,
	"151ace02-771b-44d1-afea-0ca3056c3b2c": 1,
	"59f2b03d-40c1-405b-aeea-38df617b49fa": 1,
	"7a1b419a-fbd7-4507-93e3-d6dba20ea7ad": 1,
	"4da81432-6ff0-4ed4-9df0-119b19ea8681": 1,
	"291a5adb-cc71-4eda-be0d-5a0af6bea973": 1,
	// 1x18
	"8f3296c5-9056-4fe5-b9c6-d0d0db113080": 1,
	// 2x10
	"63387c79-e174-4d78-9d5e-86fdaa8bc45e": 2,
	// 2x12
	"2af64b91-de12-48cf-a90e-2e5c4bcc9530": 2,
	"1d1c171f-7380-4727-8bd6-c63021db9dd6": 2,
	"8e98245b-fd88-48b4-8951-2d6171c717ce": 2,
	"3c9256bb-9b75-4dd2-8d49-c9236eedfcb9": 2,
	"f7902634-12e9-4a2d-9f9a-bcd22781cdab": 2,
	"3e1e42ab-294b-4816-953e-ef800bfdbbe6": 2,
	"afe74c01-d70d-47ba-bfd5-4ac43ecaa954": 2,
	"d9083a81-5a8b-4b19-990e-4d933bed5067": 2,
	"16b2e136-c230-4598-9931-5b913efb76e6": 2,
	"f5e43052-b042-4f7c-b49d-a96c6571a4ce": 2,
	"d004a94e-14b6-4e92-80eb-89aa306f0008": 2,
	"356da432-3493-4403-ac2d-761a28855638": 2,
	"b1de33d6-c660-4bf6-b290-79626a35149e": 2,
	"bce514e5-5fa0-40c3-ae1d-85093a3d62f6": 2,
	"8b37839c-5798-4584-8aca-b7bfce4819a6": 2,
	"a56735e4-226b-4887-bbb7-9a04fa080c0f": 2,
	"6ac86da3-7341-4a2f-b9c4-ecd118c559c4": 2,
	// 2x15
	"17d311d6-440e-4d5b-a695-f74d55e986dd": 2,
	"a55c9112-55bf-4fc2-845d-fd4ce59eb40c": 2,
	// 3x10
	"21a4eb35-fce4-482a-9199-ace3f4205ede": 3,
	"104346d1-5fbc-474d-bcf0-7980cbe04a82": 3,
	// 4x10
	"fa7d0a6a-23e9-4f05-a062-607f5cfe7dda": 4,
	"c4af1d27-afd0-426f-9c96-e2f65378b93a": 4,
	"9fa8c924-6543-4085-b55b-58b99aada17e": 4,
	"6ea64fbe-8fd5-4eac-a216-28cb0f07faa3": 4,
	"3b999e93-6120-4a91-b3d7-0f2902be747f": 4,
	"4614704e-7ca2-4736-a750-648fe9033650": 4,
	"fad46e0e-3760-484f-9f7b-85f0376780ef": 4,
	"1052d32a-7bd7-4463-98c0-b9cc18a000be": 4,
	// 4x12
	"936efc52-2172-4faf-9be5-e8f45244a2b9": 4,
	"c97bc69c-c02d-4cce-b19d-859b72833550": 4,
	"cd231131-d053-4193-bd88-7bbd14144680": 4,
	"7c0b8ce1-cbb4-4e5b-9973-a572143ddb2b": 4,
	"a54d97da-cd7f-4742-acba-45cb2688c8c9": 4,
	"c6dc5147-0436-482f-9a8d-070ccea23c46": 4,
	"79818fe5-e89f-437b-86bc-de80ecea216d": 4,
	"c4ea21cc-6444-4779-9eee-62d4bc085410": 4,
	"9ef10de4-2781-4ab7-9179-c1cc4b6615e5": 4,
	"0f086cfd-b793-4b69-894a-b69d4c32154b": 4,
	"67f95a0d-34e8-4206-b321-3e57c8d1b407": 4,
	"445c0a64-c729-4502-b7c4-91e211a7fc21": 4,
	"8a9bd7a7-d080-4023-9a60-e8cfcddccf0e": 4,
	"8b147712-d44d-4564-a40f-fe288110ea6c": 4,
	"2bca68bf-1cbd-4b37-b6c2-1b40e1902e5c": 4,
	"ff9cad13-059d-46bf-815a-d8f851c410d4": 4,
	"4edd00f5-1dc0-4130-a3bb-7eb608e834ab": 4,
	"849b3340-9e28-411f-9faf-e99b7b2bfb36": 4,
	"6dfb576f-b549-4dd4-ac79-1021bcb53bb2": 4,
	"81866d8a-1072-4bcd-adce-861f842f4e5c": 4,
	"8eaae0e6-4c9b-471e-814b-6a4596a4d153": 4,
	"088fe7ff-34f6-4a64-9a7c-0f938f05553d": 4,
	// 8x10
	"fb54d283-dfc2-402e-8c76-87713387d770": 8,
	"b3900806-4600-4cc0-a2f2-217b9ec09c5f": 8,
	"b31f4357-3c42-4c89-8736-64f25bcbed9d": 8,
	"b11418d6-9a01-42a9-a85b-cd287ac7a98e": 8,
	// Other
	"846d771b-3037-4d47-90c7-30e33e32af6f": 4,
	"c275ee33-8180-4b0c-9215-dd6f37aa2394": 4,
}

var Speakers = map[string]string{
	"a3cc18b8e9b449e3b1ce34c69b310b83": "American 12C",
	"d2b5f9c3e33d442ab14cf65d84aed0f5": "American 12K",
	"02079eab6ff44741961cd95bb82b9662": "American Alnico",
	"0b4e1019fe2d42c7b2292a29d4194543": "American Bulldog",
	"e372dd04b11d49588c290fbe341e97ca": "Brit 75",
	"942153d281fb4b089fc20e07a34e9ca7": "Brit 80",
	"d9c445a5002341f191b0c066d4a45eb3": "Brit 100",
	"aa7f635a7c284116a6229675340f9fd8": "Brit Alnico B",
	"96d52a2264b8495bb0a5c2571deb498f": "Brit Alnico G",
	"492ec44546cb43798742ddd231cf632a": "Brit Alnico S",
	"674b563d948e4f3398d18f8904096315": "Brit Anniversary 1",
	"7f26988d1b424e83b12587238b83c623": "Brit Anniversary 2",
	"7b2ac1f3a2f1478babce98766d5e2cd8": "Brit Darkness",
	"a56188a9a6bc4373903dbbde779548f1": "Brit Green",
	"93ece316161d4a7db5c075a64a873b02": "Brit Silver",
	"fc5bcd9eedca47b786e18803eb284b9c": "Brit T12G",
	"5eca5662178a403885f7caea56cef141": "Brit V1",
	"2dc1a3c46a204deba9cd5e939ae1e1fa": "Brit V2",
	"8c9127bce65e47f1bfd6c873fdbe822d": "Brit V3",
	"1a8ca2dad6434218b82dbf98921c0a9b": "Brit Vintage 8",
	"b413c57dca9541778646330ee16375c5": "Brit Vintage 16A",
	"9422a3d95e6b4c63bc6db15fcbd99f09": "Brit VIntage 16B",
	"8b9fc1cef0124429b728f0e822a1329e": "California Red",
	"91e0a91609a74704b7739023e3bda8d8": "Custom Fender",
	"d61186a940d948d1889194a5e6dcfc6b": "CV GT12-16",
	"4c176b93da64461bb9894d042c5475fc": "EV Darkness",
	"d052f84ca5fd4a699d4bd4fa68c155f2": "HiAmp",
	"f755dce5b3004aae8b07adac9da35705": "Jazz 12",
	"a5cad4f1d3b144ceaf6f6b4c2094cddc": "Metal V 1200",
	"a13a9305422c4f2f893fa58dcdca4f2e": "Silver Alnico",
}

var Mics = map[string]string{
	// Condensers
	"1425abc1-2525-4d85-bfbf-f40009c2f19c": "Bottle 563",
	"035c9475-312b-4f6e-87b4-c33aad7d5470": "Condenser 12",
	"9f8a2c8c-aa21-43ab-a316-5084479de02e": "Condenser 67",
	"2b667232-0a83-4132-a18b-f51e71fa349c": "Condenser 84",
	"9e444286-cab4-46a4-bfa3-a6d55b3ffcfb": "Condenser 87",
	"d78598e7-3fa4-46da-aadc-5a7733b7f896": "Condenser 170",
	"0f35a776-f6db-403d-930f-6b7f42fed749": "Condenser 414",
	"8e0525ab-e522-41d3-870e-9da851c42167": "MD1-b",
	"333890d1-62de-4f2a-a4c3-1ca0dd0d9196": "Tube VM",
	// Dynamic
	"eb1d233a-8aec-4708-b42b-b4fd26397889": "Dynamic 20",
	"1e41acc4-85af-4e84-bee4-eabc0be5fef1": "Dynamic 57",
	"b216abec-6fae-4fcd-95fd-c89aacf60ee2": "Dynamic 421",
	"c8fce7b6-deaf-461c-8628-5cfe82c15173": "Dynamic 441",
	"373859a6-cfc1-4c2c-ab8c-35ddbfb8ee77": "Dynamic 609",
	"565a6dcf-89df-4190-a552-c76d78bdab66": "Vintage Dynamic 20",
	// Overheads
	"Condenser 12":  "Condenser 12",
	"Condenser 87":  "Condenser 87",
	"Condenser 170": "Condenser 170",
	"Condenser 414": "Condenser 414",
	// Ribbon
	"cf06582b-4b26-42ce-9491-e00e7ab2481e": "Ribbon 121",
	"f1869200-4515-4ab5-a690-096d142e548d": "Ribbon 160",
	"1cb1f17b-bc70-485d-bd8a-339e54eedac5": "Velo-8",
}

var Rooms = map[string]string{
	"Amp Closet":   "Amp Closet",
	"Bathroom":     "Bathroom",
	"Garage":       "Garage",
	"Hall":         "Hall",
	"Large Studio": "Large Studio",
	"Mid Studio":   "Mid Studio",
	"Small Studio": "Small Studio",
	"Subway":       "Subway",
}

var FX = map[string]string{
	// Pedals
	// Delay
	"ad9d0a70-7a59-4fef-ace5-c592764e3749": "'63 Reverb",
	"b756e0c1-7685-4b38-bccc-b74c7febd868": "Analog Delay",
	"e11b1dc5-1f7d-42ad-af30-0539b3646b3c": "Delay",
	"48e7b721-d57a-4c34-813b-95d8091d5eda": "EchoMan",
	"907ecdf1-15be-4f41-b56d-2705e7bb89ae": "EP Tape Echo",
	"bf72ebc2-a539-4cd2-9204-2d91e9d573df": "Replica",
	"4468f4f7-0068-4b8b-ac2b-99e13113fe2d": "Slash Delay",
	"28bb2c33-0bdf-44f7-9274-2eca934cbbff": "SSTE",
	"96b57f95-4380-444a-8c0a-fbcc9bef1dd9": "TapDelay",
	"8bbfc5b9-bf29-4a55-8211-ca21dcfda8bf": "Tape Echo",
	// Distortion
	"58dbec22-58e0-464c-8c04-91fb9d9973e2": "BigPig",
	"305c9b6b-04cf-4673-b58a-e62afb4fefcb": "Crusher",
	"5e65abef-82eb-4995-b911-d5eca4f8291e": "Diode Overdrive",
	"510f6d25-6ec4-417b-bf58-0f8028209cce": "Distortion",
	"395ed825-f3e8-40c1-8d69-34d8b23c9100": "Feedback",
	"e5c8acd3-3771-4df9-8d2e-ee33c8dd3d21": "Metal Distortion 2",
	"1910832b-2b47-46ff-b14c-46ec168e50e6": "Metal Distortion",
	"1d03a910-c5a3-461e-a43a-485ddf3d84ef": "Moller",
	"e5644c95-e382-4cfe-9c1f-85451017771d": "Mudhoney",
	"7c499158-084f-49b1-9543-f7e9acc122e0": "OCD",
	"967e57ac-b67d-4b97-942e-aca407e306e0": "OctoBlue",
	"fd627f5e-ba11-4082-b546-a4f0b05985ff": "Overdrive",
	"fa1de2e2-102b-4edf-b3b5-23ceaeddedf0": "Overscream",
	"8a96f6a6-49af-41fb-ab36-a62a18f17def": "Pinnacle Deluxe",
	"16daf2e6-1c56-4abe-97c9-1fffe2b22bb2": "Power Grid",
	"9b672f82-2832-4134-8db7-5cb9147c69a3": "PRODrive",
	"1d665fde-1a62-42a1-be6d-bad9bbe5df3d": "SVX-OD",
	"c8b142b0-4480-4d79-bc5c-f0232440ce05": "The Ambass'dor",
	// Dynamics
	"77f0f320-cc4e-44be-9ffe-2f0b679434ae": "Booster",
	"5478981b-b18a-469f-81e7-a3e228cc9d50": "Compressor",
	"26c75920-d4bf-4e5e-900f-f78c70e06c17": "Dcomp",
	"f5edced9-6dfc-4851-8651-f81f5423d210": "Fender Compressor",
	"d3e05ec0-2c7b-498a-adc0-b263e853ad30": "Gate",
	"0455f997-43ca-4c9b-9269-286a19d10d48": "Noise Gate",
	"8a24aa96-f0ae-4e1c-a534-6671e245a690": "SVX Compressor",
	// EQ
	"8d7ff76e-9273-46b6-95d5-3d7bd667fff2": "7 Band Graphic",
	"babadeaf-9c28-4641-8fa9-d7366a3238a2": "10 Band Graphic",
	// Filter
	"15b140e0-3e02-4adc-a9c4-c652960e60f9": "Bass Envelope Filter",
	"01cadfae-3ced-4ea6-8676-29a7e6c920b2": "Bass Wah",
	"487cd1a4-834e-45b2-b5be-6a424cc6a123": "Contour Wah",
	"77a321dd-69e1-4474-be07-d8a97e78bd1f": "Envelope Filter",
	"75f96017-8a09-41fd-9979-75bf8bf81645": "Fender Wah",
	"a58d91b0-d7c5-4d3d-8a9a-5c8b75335502": "Fuzz Wah",
	"390c602d-5834-417d-bf0c-cafe544c5869": "LFO Filter",
	"0332d916-2ab2-4b7d-98c4-73a80a42b3b1": "Nu-Tron III",
	"327d6d53-b6cb-4d33-bdaf-620fb52c20ec": "Rezo",
	"25425c78-31db-48f4-ad57-09f41e0e1291": "Step Filter",
	"2de5239a-78d6-4a01-82e6-2ea3afb60501": "Wah 10",
	"bc86a019-ffd5-4b71-8bfe-5913e3d58d7c": "Wah 46",
	"6482748e-9382-4ad6-b284-5c29ee50f2d7": "Wah",
	"88863a3a-cfe3-4e86-b735-1303c511bf5f": "WahDist",
	// Fuzz
	"8beec4ce-fb43-4f81-935a-3b5cb3695c8b": "Class Fuzz",
	"09ac5b94-f238-4e4c-914e-ba7662f280d9": "Fuzz Age 2",
	"6c3ff0bf-b840-47f3-83d3-66816763097f": "Fuzz Age",
	"0679dea3-2588-4d9d-8d0d-ef3762f1f478": "Fuzz One",
	"aa74a915-a1fe-4f54-a8a8-5297c3e09b56": "Octa-V",
	"b0f5949f-4825-4202-92a0-c5817f493116": "RightFuzz",
	"64e7c1cd-b860-40c7-930b-6d820b1ffa77": "XS Fuzz",
	// Modulation
	"ae6177c2-27c2-4463-a06a-357408bb2082": "Analog Flanger",
	"ed2c3a06-d304-496b-b031-7725a3d27eea": "Bass Analog Chorus",
	"bc6a9f33-ac11-41f8-973d-0327d4f3e018": "Chorus",
	"2a9ef349-fb29-4e66-99a9-cc66d10192cc": "Chorus-1",
	"8a878202-9126-4d20-8e73-374e178312f4": "Electric Flanger",
	"7ccf016f-e540-4e46-a124-8f19ce5ab2b1": "Flanger",
	"4e4d82f9-224a-4ffb-9994-97ef8285c315": "Metal Flanger",
	"b1ad4a5d-1ad2-4b32-8532-945b869409e3": "Nirvana",
	"50378f09-a919-4dee-9bbe-c242403a52a2": "Opto Tremolo",
	"6178531f-d021-43c0-8922-858ffa085746": "Phaser",
	"a4ed5e25-707d-40ef-9846-64eeb820aeea": "Phaze Nine",
	"cc424097-15e5-47d3-abb9-3925073ac22b": "Phazer10",
	"86875e91-6fbd-4198-a45c-a06119e6a967": "Seek Trem",
	"0ba47121-179c-4d42-bbb6-c3e81bb4f7af": "Seek Wah",
	"96ae9a18-1c2b-48cc-843a-851adb43c091": "Shape Shifter",
	"0ef53d8f-2dd5-4acd-95f8-e8652ae31240": "Small Phazer",
	"187eb9ab-7ae6-4797-954b-079de09e26bb": "Tremolo",
	"a6d48956-a0e5-4d63-9c22-b5b38604d2a5": "Uni-V",
	"5f3947b1-6a09-4570-9f9c-1cc53a7fd88f": "X-Chorus",
	// Other
	"71fe6e6d-5879-42a7-9a31-6093ecee2a1c": "Acoustic Sim",
	"01776ae8-8442-4633-b5f7-6bfdaf423ccb": "Fender Volume",
	"66410529-1158-4d6e-a33a-474541a64571": "Step Slicer",
	"7b1dc197-a4ac-41cc-8b1e-d8ed4102f432": "SVX Volume",
	"ca453f6e-7af5-4e90-90df-ff954b17ecc2": "Swell",
	"de12969a-31cc-4985-b4cf-289d2970823d": "Volume",
	// Pitch
	"01648ef1-6369-4170-81a3-90dd20451260": "Blender",
	"46f09ab5-ffd9-4c5b-8eec-681f880d4530": "Harmonator",
	"994770ae-ebb4-4ca8-884e-374f88fa3db0": "Octav",
	"e2b29e5c-33a0-41f0-9d54-dc749d371fe0": "Pitch Shifter",
	"9afc331b-c0c3-4592-b03f-c97f8d911e34": "SVX-OCT",
	"9b8e89e2-2959-41b2-90eb-dc5de12964d0": "Wharmonator",
	// Rack Effects
	// Delay
	"1189979a-db5d-4dc1-9228-7bd974d8a8c5": "Digital Delay",
	"773b8ea7-b54a-4a3c-99df-ffbbf6d29271": "Tap Delay",
	"a8a839aa-35e1-4fac-8834-a0a1701c63d8": "Tape Echo",
	// Dynamics
	"7307c816-856f-438b-a381-45edf43bee0b": "Compressor",
	"d0211742-18e6-4fdb-9efa-3d72e4ae515b": "Tube Compressor",
	// EQ
	"b66b51c2-d9a3-4909-b7e0-cd1e51636e97": "Graphic EQ",
	"9f1147a6-302f-48f3-a5bc-26cc5d399a8b": "Parametric EQ 3",
	// Filter
	"e2e5495c-5ac3-405f-9fdd-b73670d413c0": "Rezo",
	"fb5d2469-05f6-4a44-9576-41ae232c9385": "Step Filter",
	"5a6dfdc0-69d2-4e84-a84c-e500a0d75505": "Wah",
	// Modulation
	"02643125-de84-4c94-b214-4d300652332b": "Analog Chorus",
	"1edbb450-d048-11dc-95ff-0800200c9a66": "Digital Chorus",
	"c11388bb-6326-4766-a440-ea9fa3f82425": "Digital Flanger",
	"99c5d753-57e3-40a4-9612-04623ac61289": "Rotary Speaker",
	"9fa5b238-d7d0-47ac-a2e3-6e4e11761261": "Sine Flange",
	"4b91de5f-73c6-46d2-957b-6b9451abf050": "Step Slicer",
	"fe891a4f-6098-423d-b8dd-3213373b990c": "Stereo Enhancer",
	"1e27e673-20fe-474e-a438-d85a9bc566b4": "Swell",
	"14fd2d3b-a81d-4850-a2a6-9e94b7351059": "TERC",
	"cee174c4-821c-4b92-8cb4-86c38c433668": "Triangle Chorus",
	"647b8569-e3b4-48c3-b8a1-37c5f920e3f6": "Harmonator",
	"0f304b4d-65b9-4347-9f44-fcaa8509efaf": "Pitch Shift",
	"845b672b-255f-4edf-9e67-68b607dcf63a": "Pitch Shifter",
	"3c8d23d7-959a-4479-b9c2-46af9a77ba46": "'63 Reverb",
	"59ab0817-b168-4bdc-b837-e3cba1efb2dd": "Digital Reverb",
}
