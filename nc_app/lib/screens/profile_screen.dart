import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';
import '../services/api_service.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({super.key});

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  Map<String, dynamic>? _profile;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final p = await ApiService.get('/api/profile');
      setState(() {
        _profile = p != null ? Map<String, dynamic>.from(p as Map) : null;
        _error = null;
      });
    } catch (e) {
      setState(() {
        _profile = null;
        _error = e.toString();
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final auth = context.watch<AuthProvider>();
    if (_error != null) {
      return Center(child: Column(mainAxisSize: MainAxisSize.min, children: [
        Text(_error!, style: const TextStyle(color: Colors.red)),
        TextButton(onPressed: _load, child: const Text('Повторить')),
      ]));
    }
    final p = _profile ?? {};
    final phone = p['phone']?.toString() ?? '';
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Имя: ${p['full_name'] ?? auth.fullName ?? ""}', style: Theme.of(context).textTheme.titleMedium),
                Text('Email: ${p['email'] ?? auth.email ?? ""}'),
                Text('Телефон: ${phone.isNotEmpty ? phone : "—"}'),
                Text('Роль: ${p['role'] ?? auth.role ?? ""}'),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: _ChangePasswordForm(initialPhone: phone),
          ),
        ),
      ],
    );
  }
}

class _ChangePasswordForm extends StatefulWidget {
  final String initialPhone;
  const _ChangePasswordForm({this.initialPhone = ''});
  @override
  State<_ChangePasswordForm> createState() => _ChangePasswordFormState();
}

class _ChangePasswordFormState extends State<_ChangePasswordForm> {
  late final TextEditingController _password;
  late final TextEditingController _phone;

  @override
  void initState() {
    super.initState();
    _password = TextEditingController();
    _phone = TextEditingController(text: widget.initialPhone);
  }
  bool _loading = false;

  @override
  void dispose() {
    _password.dispose();
    _phone.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final body = <String, dynamic>{};
    if (_password.text.isNotEmpty) body['password'] = _password.text;
    if (_phone.text.isNotEmpty) body['phone'] = _phone.text;
    if (body.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Введите пароль и/или телефон')));
      return;
    }
    setState(() => _loading = true);
    try {
      await ApiService.put('/api/profile', body);
      if (mounted) {
        _password.clear();
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Сохранено')));
      }
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        const Text('Сменить пароль / телефон', style: TextStyle(fontWeight: FontWeight.bold)),
        const SizedBox(height: 12),
        TextField(controller: _password, decoration: const InputDecoration(labelText: 'Новый пароль', border: OutlineInputBorder()), obscureText: true),
        const SizedBox(height: 8),
        TextField(controller: _phone, decoration: const InputDecoration(labelText: 'Телефон', border: OutlineInputBorder()), keyboardType: TextInputType.phone),
        const SizedBox(height: 12),
        FilledButton(onPressed: _loading ? null : _submit, child: _loading ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2)) : const Text('Сохранить')),
      ],
    );
  }
}
