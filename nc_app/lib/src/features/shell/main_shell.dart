import 'package:flutter/material.dart';
import '../../services/auth_state.dart';
import '../home/home_screen.dart';

class MainShell extends StatefulWidget {
  const MainShell({super.key, required this.authState});

  final AuthState authState;

  @override
  State<MainShell> createState() => _MainShellState();
}

class _MainShellState extends State<MainShell> {
  int _selectedIndex = 0;

  static const _studentTeacherItems = [
    _NavItem(Icons.home_outlined, Icons.home, 'Главная'),
    _NavItem(Icons.calendar_today_outlined, Icons.calendar_today, 'Учебный процесс'),
    _NavItem(Icons.people_outline, Icons.people, 'Социум'),
    _NavItem(Icons.person_outline, Icons.person, 'Профиль'),
  ];

  List<_NavItem> get _navItems {
    final role = widget.authState.role;
    if (role == 'admin' || role == 'director') {
      return [
        ..._studentTeacherItems,
        _NavItem(Icons.admin_panel_settings_outlined, Icons.admin_panel_settings, 'Админ'),
      ];
    }
    return _studentTeacherItems;
  }

  Widget _bodyForIndex(int index) {
    switch (index) {
      case 0:
        return HomeScreen(
          apiToken: widget.authState.token!,
          onLogout: () async {
            await widget.authState.clear();
          },
        );
      case 1:
        return const _PlaceholderScreen(title: 'Учебный процесс', subtitle: 'Расписание, журнал, файлы');
      case 2:
        return const _PlaceholderScreen(title: 'Социум', subtitle: 'Форум, чат');
      case 3:
        return _ProfilePlaceholder(authState: widget.authState);
      case 4:
        return const _PlaceholderScreen(title: 'Админ', subtitle: 'Управление контентом и пользователями');
      default:
        return HomeScreen(
          apiToken: widget.authState.token!,
          onLogout: () async => await widget.authState.clear(),
        );
    }
  }

  @override
  Widget build(BuildContext context) {
    final isNarrow = MediaQuery.sizeOf(context).width < 700;
    return Scaffold(
      body: Row(
        children: [
          if (!isNarrow)
            NavigationRail(
              selectedIndex: _selectedIndex,
              onDestinationSelected: (i) => setState(() => _selectedIndex = i),
              labelType: NavigationRailLabelType.all,
              destinations: _navItems
                  .map((e) => NavigationRailDestination(
                        icon: Icon(e.outlined),
                        selectedIcon: Icon(e.filled),
                        label: Text(e.label),
                      ))
                  .toList(),
            ),
          Expanded(
            child: _bodyForIndex(_selectedIndex),
          ),
        ],
      ),
      bottomNavigationBar: isNarrow
          ? NavigationBar(
              selectedIndex: _selectedIndex,
              onDestinationSelected: (i) => setState(() => _selectedIndex = i),
              destinations: _navItems
                  .map((e) => NavigationDestination(
                        icon: Icon(e.outlined),
                        selectedIcon: Icon(e.filled),
                        label: e.label,
                      ))
                  .toList(),
            )
          : null,
    );
  }
}

class _NavItem {
  const _NavItem(this.outlined, this.filled, this.label);
  final IconData outlined;
  final IconData filled;
  final String label;
}

class _PlaceholderScreen extends StatelessWidget {
  const _PlaceholderScreen({required this.title, this.subtitle});

  final String title;
  final String? subtitle;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(title, style: Theme.of(context).textTheme.headlineSmall),
          if (subtitle != null) ...[
            const SizedBox(height: 8),
            Text(subtitle!, style: Theme.of(context).textTheme.bodyMedium?.copyWith(color: Colors.grey)),
          ],
        ],
      ),
    );
  }
}

class _ProfilePlaceholder extends StatelessWidget {
  const _ProfilePlaceholder({required this.authState});

  final AuthState authState;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(authState.fullName ?? authState.email ?? 'Профиль', style: Theme.of(context).textTheme.headlineSmall),
            const SizedBox(height: 8),
            Text(authState.email ?? '', style: Theme.of(context).textTheme.bodyMedium?.copyWith(color: Colors.grey)),
            const SizedBox(height: 24),
            OutlinedButton(
              onPressed: () async {
                await authState.clear();
              },
              child: const Text('Выйти'),
            ),
          ],
        ),
      ),
    );
  }
}
