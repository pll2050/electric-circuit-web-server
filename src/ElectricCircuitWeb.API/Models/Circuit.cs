namespace ElectricCircuitWeb.API.Models;

/// <summary>
/// 회로 모델
/// </summary>
public class Circuit
{
    public string Id { get; set; } = Guid.NewGuid().ToString();
    public string ProjectId { get; set; } = string.Empty;
    public string Name { get; set; } = string.Empty;
    public string Data { get; set; } = string.Empty; // JSON으로 저장되는 회로 데이터
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    public DateTime? UpdatedAt { get; set; }
}
