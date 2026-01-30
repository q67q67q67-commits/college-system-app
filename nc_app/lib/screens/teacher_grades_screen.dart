import 'package:flutter/material.dart';
import '../services/api_service.dart';

class TeacherGradesScreen extends StatefulWidget {
  const TeacherGradesScreen({super.key});

  @override
  State<TeacherGradesScreen> createState() => _TeacherGradesScreenState();
}

class _TeacherGradesScreenState extends State<TeacherGradesScreen> {
  List<dynamic> _groups = [];
  List<dynamic> _schedule = [];
  List<dynamic> _students = [];
  int? _selectedGroupId;
  int? _selectedScheduleId;
  int? _selectedStudentId;
  final _grade = TextEditingController();
  final _comment = TextEditingController();
  DateTime _date = DateTime.now();
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  @override
  void dispose() {
    _grade.dispose();
    _comment.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    try {
      final g = await ApiService.get('/api/groups');
      final s = await ApiService.get('/api/schedule');
      setState(() {
        _groups = g is List ? g : [];
        _schedule = s is List ? s : [];
        _error = null;
      });
    } catch (e) {
      setState(() {
        _groups = _schedule = [];
        _error = e.toString();
      });
    }
  }

  Future<void> _loadStudents(int groupId) async {
    try {
      final res = await ApiService.get('/api/groups/$groupId/students');
      setState(() => _students = res is List ? res : []);
    } catch (_) {
      setState(() => _students = []);
    }
  }

  Future<void> _submit() async {
    if (_selectedStudentId == null || _selectedScheduleId == null || _grade.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Выберите студента, занятие и введите оценку')));
      return;
    }
    try {
      await ApiService.post('/api/grades', {
        'user_id': _selectedStudentId,
        'schedule_id': _selectedScheduleId,
        'grade': double.tryParse(_grade.text.trim()) ?? 0,
        'grade_date': '${_date.year}-${_date.month.toString().padLeft(2, '0')}-${_date.day.toString().padLeft(2, '0')}',
        'comment': _comment.text.trim(),
      });
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Оценка выставлена')));
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_error != null) return Center(child: Column(mainAxisSize: MainAxisSize.min, children: [Text(_error!, style: const TextStyle(color: Colors.red)), TextButton(onPressed: _load, child: const Text('Повторить'))]));
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        DropdownButtonFormField<int>(
          value: _selectedGroupId,
          decoration: const InputDecoration(labelText: 'Группа', border: OutlineInputBorder()),
          items: _groups.map((g) => DropdownMenuItem<int>(value: (g['id'] as num).toInt(), child: Text(g['name']?.toString() ?? ''))).toList(),
          onChanged: (v) {
            setState(() {
              _selectedGroupId = v;
              _selectedStudentId = null;
              _students = [];
            });
            if (v != null) _loadStudents(v);
          },
        ),
        const SizedBox(height: 8),
        DropdownButtonFormField<int>(
          value: _selectedStudentId,
          decoration: const InputDecoration(labelText: 'Студент', border: OutlineInputBorder()),
          items: _students.map((s) => DropdownMenuItem<int>(value: (s['user_id'] as num).toInt(), child: Text(s['full_name']?.toString() ?? ''))).toList(),
          onChanged: (v) => setState(() => _selectedStudentId = v),
        ),
        const SizedBox(height: 8),
        DropdownButtonFormField<int>(
          value: _selectedScheduleId,
          decoration: const InputDecoration(labelText: 'Занятие', border: OutlineInputBorder()),
          items: _schedule.map((s) => DropdownMenuItem<int>(value: (s['id'] as num).toInt(), child: Text('${s['subject']} · ${s['room']}' ?? ''))).toList(),
          onChanged: (v) => setState(() => _selectedScheduleId = v),
        ),
        const SizedBox(height: 8),
        TextField(controller: _grade, decoration: const InputDecoration(labelText: 'Оценка', border: OutlineInputBorder()), keyboardType: TextInputType.number),
        const SizedBox(height: 8),
        ListTile(title: const Text('Дата'), subtitle: Text('$_date'), onTap: () async { final d = await showDatePicker(context: context, initialDate: _date, firstDate: DateTime(2020), lastDate: DateTime(2030)); if (d != null) setState(() => _date = d); }),
        const SizedBox(height: 8),
        TextField(controller: _comment, decoration: const InputDecoration(labelText: 'Комментарий', border: OutlineInputBorder()), maxLines: 2),
        const SizedBox(height: 16),
        FilledButton(onPressed: _submit, child: const Text('Выставить оценку')),
      ],
    );
  }
}
