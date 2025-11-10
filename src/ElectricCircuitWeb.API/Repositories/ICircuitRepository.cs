using ElectricCircuitWeb.API.Models;

namespace ElectricCircuitWeb.API.Repositories;

/// <summary>
/// 회로 리포지토리 인터페이스
/// </summary>
public interface ICircuitRepository
{
    Task<IEnumerable<Circuit>> GetByProjectIdAsync(string projectId);
    Task<Circuit?> GetByIdAsync(string circuitId);
    Task<Circuit> CreateAsync(Circuit circuit);
    Task<Circuit> UpdateAsync(Circuit circuit);
    Task<bool> DeleteAsync(string circuitId);
}
