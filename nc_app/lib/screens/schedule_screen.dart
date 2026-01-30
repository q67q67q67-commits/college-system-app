import 'package:flutter/material.dart';
import '../services/api_service.dart';

class ScheduleScreen extends StatefulWidget {
  const ScheduleScreen({super.key});

  @override
  State<ScheduleScreen> createState() => _ScheduleScreenState();
}

class _ScheduleScreenState extends State<ScheduleScreen> {
  List<dynamic> _list = [];
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final res = await ApiService.get('/api/schedule');
      setState(() {
        _list = res is List ? res : [];
        _error = null;
      });
    } catch (e) {
      setState(() {
        _list = [];
        _error = e.toString();
      });
    }
  }

  static const _days = ['', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'];

  @override
  Widget build(BuildContext context) {
    if (_error != null) return Center(child: Column(mainAxisSize: MainAxisSize.min, children: [Text(_error!, style: const TextStyle(color: Colors.red)), TextButton(onPressed: _load, child: const Text('Повторить'))]));
    final byDay = <int, List<dynamic>>{};
    for (final s in _list) {
      final d = (s['day_of_week'] as num?)?.toInt() ?? 0;
      byDay.putIfAbsent(d, () => []).add(s);
    }
    return RefreshIndicator(
      onRefresh: _load,
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [1, 2, 3, 4, 5, 6, 7].map((d) {
          final items = byDay[d] ?? [];
          if (items.isEmpty) return const SizedBox.shrink();
          return Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(_days[d], style: Theme.of(context).textTheme.titleMedium),
              ...items.map<Widget>((s) => Card(
                    margin: const EdgeInsets.only(bottom: 4),
                    child: ListTile(
                      title: Text(s['subject']?.toString() ?? ''),
                      subtitle: Text('${s['start_time']}–${s['end_time']} · Ауд. ${s['room'] ?? '—'}'),
                    ),
                  )),
              const SizedBox(height: 12),
            ],
          );
        }).toList(),
      ),
    );
  }
}
