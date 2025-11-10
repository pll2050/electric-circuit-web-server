using ElectricCircuitWeb.API.Models;

namespace ElectricCircuitWeb.API.Services;

/// <summary>
/// 회로 서비스 인터페이스
/// </summary>
public interface ICircuitService
{
    // 회로 조회
    Task<IEnumerable<Circuit>> GetProjectCircuitsAsync(string projectId);
    Task<Circuit?> GetCircuitByIdAsync(string circuitId);

    // 회로 생성
    Task<Circuit> CreateCircuitAsync(Circuit circuit);

    // 회로 수정
    Task<Circuit> UpdateCircuitAsync(Circuit circuit);

    // 회로 삭제
    Task<bool> DeleteCircuitAsync(string circuitId);

    // 템플릿
    Task<IEnumerable<object>> GetTemplatesAsync();
    Task<Circuit> CreateFromTemplateAsync(string templateId, string projectId, string name);
}
