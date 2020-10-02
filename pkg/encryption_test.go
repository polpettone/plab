package pkg

import "testing"

func Test_should_encrypt_and_decrypt(t *testing.T) {

	textToEncrypt := "encrypt me"
	secretWithCorrectLenght := "1234567890123456"

	encrypted, err := encrypt(textToEncrypt, secretWithCorrectLenght)

	if err != nil {
		t.Errorf("No error expected but got %v", err)
	}

	decrypted := decrypt(encrypted, secretWithCorrectLenght)
	if decrypted != textToEncrypt {
			t.Errorf("wanted %s got %s", textToEncrypt, decrypted)
	}
}

func Test_should_return_error_when_key_length_is_not_correct(t *testing.T) {

	textToEncrypt := "encrypt me"
	secretWithWrongLenght := "1234"
	encrypted, err := encrypt(textToEncrypt, secretWithWrongLenght)

	if err == nil {
		t.Error("Error expected")
	}

	if encrypted != "" {
		t.Errorf("wanted '' got %s", encrypted)
	}


}
