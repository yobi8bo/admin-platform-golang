package crypto

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 生成密码哈希，调用方只应保存返回的哈希值。
func HashPassword(raw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword 校验明文密码是否匹配已保存的 bcrypt 哈希。
func CheckPassword(hash, raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)) == nil
}
