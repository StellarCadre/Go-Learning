// 创建时间：2026/8/18 下午6:12
package vo

//只返回前端需要的字段，隐藏删除时间、更新时间等无关字段；如果后续有密码等敏感字段，绝对不能写在 VO 里。
import "time"

// UserInfoVO 用户信息-返回给前端的视图对象
type UserInfoVO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}
