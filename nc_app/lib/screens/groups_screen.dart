import 'package:flutter/material.dart';
import '../services/api_service.dart';

class GroupsScreen extends StatefulWidget {
  const GroupsScreen({super.key});

  @override
  State<GroupsScreen> createState() => _GroupsScreenState();
}

class _GroupsScreenState extends State<GroupsScreen> {
  List<dynamic> _list = [];
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final res = await ApiService.get('/api/groups');
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
    if (_error != null) return Center(child: Column(mainAxisSize: MainAxisSize.min, children: [Text(_error!, style: const TextStyle(color: Colors.red)), TextButton(onPressed: _load, child: const Text('Повторить'))]));
    return RefreshIndicator(
      onRefresh: _load,
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: _list.length,
        itemBuilder: (_, i) {
          final g = _list[i] as Map<String, dynamic>;
          return Card(
            margin: const EdgeInsets.only(bottom: 8),
            child: ListTile(
              title: Text(g['name']?.toString() ?? ''),
              subtitle: Text('${g['description'] ?? ''} · Студентов: ${g['student_count'] ?? 0}'),
              trailing: IconButton(
                icon: const Icon(Icons.people),
                onPressed: () async {
                  try {
                    final students = await ApiService.get('/api/groups/${g['id']}/students');
                    if (!context.mounted) return;
                    showModalBottomSheet(
                      context: context,
                      builder: (_) => ListView(
                        shrinkWrap: true,
                        children: (students is List ? students : []).map<Widget>((s) => ListTile(title: Text(s['full_name']?.toString() ?? ''), subtitle: Text(s['email']?.toString() ?? ''))).toList(),
                      ),
                    );
                  } catch (_) {}
                },
              ),
            ),
          );
        },
      ),
    );
  }
}
