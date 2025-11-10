using ElectricCircuitWeb.API.Models;
using ElectricCircuitWeb.API.Repositories;

namespace ElectricCircuitWeb.API.Services;

/// <summary>
/// 회로 서비스 구현
/// </summary>
public class CircuitService : ICircuitService
{
    private readonly ICircuitRepository _circuitRepository;
    private readonly ILogger<CircuitService> _logger;

    public CircuitService(ICircuitRepository circuitRepository, ILogger<CircuitService> logger)
    {
        _circuitRepository = circuitRepository;
        _logger = logger;
    }

    public async Task<IEnumerable<Circuit>> GetProjectCircuitsAsync(string projectId)
    {
        return await _circuitRepository.GetByProjectIdAsync(projectId);
    }

    public async Task<Circuit?> GetCircuitByIdAsync(string circuitId)
    {
        return await _circuitRepository.GetByIdAsync(circuitId);
    }

    public async Task<Circuit> CreateCircuitAsync(Circuit circuit)
    {
        // ID가 없으면 생성
        if (string.IsNullOrEmpty(circuit.Id))
        {
            circuit.Id = Guid.NewGuid().ToString();
        }

        circuit.CreatedAt = DateTime.UtcNow;
        circuit.UpdatedAt = null;

        return await _circuitRepository.CreateAsync(circuit);
    }

    public async Task<Circuit> UpdateCircuitAsync(Circuit circuit)
    {
        var existingCircuit = await _circuitRepository.GetByIdAsync(circuit.Id);
        if (existingCircuit == null)
        {
            throw new KeyNotFoundException($"Circuit with ID {circuit.Id} not found");
        }

        return await _circuitRepository.UpdateAsync(circuit);
    }

    public async Task<bool> DeleteCircuitAsync(string circuitId)
    {
        return await _circuitRepository.DeleteAsync(circuitId);
    }

    public async Task<IEnumerable<object>> GetTemplatesAsync()
    {
        // TODO: 실제 템플릿 구현
        return await Task.FromResult(new List<object>());
    }

    public async Task<Circuit> CreateFromTemplateAsync(string templateId, string projectId, string name)
    {
        // TODO: 실제 템플릿 기반 회로 생성 구현
        // 현재는 빈 회로 생성
        var circuit = new Circuit
        {
            Id = Guid.NewGuid().ToString(),
            ProjectId = projectId,
            Name = name,
            Data = "{}", // 빈 JSON
            CreatedAt = DateTime.UtcNow
        };

        return await _circuitRepository.CreateAsync(circuit);
    }
}
