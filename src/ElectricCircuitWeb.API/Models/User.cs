namespace ElectricCircuitWeb.API.Models;

/// <summary>
/// 사용자 모델
/// </summary>
public class User
{
    public int Id { get; set; }
    public string FirebaseUid { get; set; } = string.Empty;
    public string Email { get; set; } = string.Empty;
    public string DisplayName { get; set; } = string.Empty;
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    public DateTime? UpdatedAt { get; set; }
}
