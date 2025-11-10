namespace ElectricCircuitWeb.API.Config;

/// <summary>
/// Firebase 설정
/// </summary>
public class FirebaseConfig
{
    public string ProjectId { get; set; } = string.Empty;
    public string ServiceAccountKeyPath { get; set; } = string.Empty;
    public string DatabaseUrl { get; set; } = string.Empty;
}
