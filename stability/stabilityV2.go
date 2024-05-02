package stability

// path: stability/stability.go

//https://platform.stability.ai/docs/api-reference

/*
Method: POST
Headers:
- authorization: Bearer ${API_KEY}
- content-type: multipart/form-data
- accept: image/jpeg, image/png, image/gif, image/webp -- OR -- application/json to receive image as base64 string

Body: form-data
Required: prompt ([1 .. 10000] characters)
Optional:
- aspect_ratio (default 1:1; enum)
- negative_prompt (< 10000 chars)
- seed (0 .. 4294967294)
- style_preset (enum)
- output_format (default jpeg; enum)

Outputs:
- byte array of the generated image
- The resolution of the generated image will be 1.5 megapixels

Credits:
Flat rate of 3 credits per successful generation. You will not be charged for failed generations.
*/
