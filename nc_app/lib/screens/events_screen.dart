import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';
import '../services/api_service.dart';

class EventsScreen extends StatefulWidget {
  const EventsScreen({super.key});

  @override
  State<EventsScreen> createState() => _EventsScreenState();
}

class _EventsScreenState extends State<EventsScreen> {
  List<dynamic> _list = [];
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final res = await ApiService.get('/api/events?limit=20');
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

  @override
  Widget build(BuildContext context) {
    final auth = context.watch<AuthProvider>();
    final canEdit = auth.isAdminOrDirector;
    if (_error != null) return Center(child: Column(mainAxisSize: MainAxisSize.min, children: [Text(_error!, style: const TextStyle(color: Colors.red)), TextButton(onPressed: _load, child: const Text('Повторить'))]));
    return RefreshIndicator(
      onRefresh: _load,
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: _list.length,
        itemBuilder: (_, i) {
          final e = _list[i] as Map<String, dynamic>;
          return Card(
            margin: const EdgeInsets.only(bottom: 8),
            child: ListTile(
              title: Text(e['title']?.toString() ?? ''),
              subtitle: Text('${e['description'] ?? ''}\n${e['event_date'] ?? ''} · ${e['location'] ?? ''}'),
              isThreeLine: true,
              trailing: canEdit
                  ? IconButton(
                      icon: const Icon(Icons.delete),
                      onPressed: () async {
                        if (!await showDialog<bool>(context: context, builder: (_) => AlertDialog(title: const Text('Удалить?'), actions: [TextButton(onPressed: () => Navigator.pop(context, false), child: const Text('Нет')), FilledButton(onPressed: () => Navigator.pop(context, true), child: const Text('Да'))])) ?? false) return;
                        try {
                          await ApiService.delete('/api/events/${e['id']}');
                          if (context.mounted) _load();
                        } catch (_) {}
                      },
                    )
                  : null,
            ),
          );
        },
      ),
    );
  }
}
