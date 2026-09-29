import { useState } from "react";
import { Plus, Pencil, Trash2, Users, Shield } from "lucide-react";
import CreateUserModal from "../modals/CreateUserModal";
import EditUserModal from "../modals/EditUserModal";

interface User {
  id: string;
  name: string;
  email: string;
  phone: string;
  role: string;
  status: "active" | "inactive";
}

const MOCK_USERS: User[] = [
  { id: "1", name: "Ahmed Al-Rashid", email: "ahmed@naseej.com", phone: "+966 55 123 4567", role: "Admin", status: "active" },
  { id: "2", name: "Sara Mohammed", email: "sara@naseej.com", phone: "+966 55 234 5678", role: "Cashier", status: "active" },
  { id: "3", name: "Omar Hassan", email: "omar@naseej.com", phone: "+966 55 345 6789", role: "Store Manager", status: "active" },
  { id: "4", name: "Layla Abdullah", email: "layla@naseej.com", phone: "+966 55 456 7890", role: "Cashier", status: "inactive" },
];

export default function UsersPage() {
  const [users] = useState<User[]>(MOCK_USERS);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  const handleEdit = (user: User) => {
    setEditingUser(user);
    setShowEditModal(true);
  };

  return (
    <>
      <div className="flex flex-col flex-1 overflow-hidden bg-background">
        {/* Header */}
        <div className="border-b border-border bg-card shadow-sm shrink-0">
          <div className="mx-auto max-w-6xl px-4 py-6 sm:px-6 lg:px-8">
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-2xl font-bold text-foreground">Users</h1>
                <p className="text-sm text-muted-foreground mt-1">
                  Manage staff accounts and access
                </p>
              </div>
              <button
                type="button"
                onClick={() => setShowCreateModal(true)}
                className="px-4 py-2 rounded-lg text-sm font-medium transition-colors bg-accent text-accent-foreground hover:bg-accent/90 flex items-center gap-2"
              >
                <Plus className="w-4 h-4" />
                Add User
              </button>
            </div>
          </div>
        </div>

        {/* Users List */}
        <div className="flex-1 overflow-y-auto mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8 w-full">
          {users.length === 0 ? (
            <div className="rounded-lg bg-card p-8 text-center">
              <Users className="w-10 h-10 text-muted-foreground/50 mx-auto mb-3" />
              <p className="text-muted-foreground">No users yet</p>
            </div>
          ) : (
            <div className="rounded-lg border border-border bg-card overflow-hidden">
              {/* Table Header */}
              <div className="grid grid-cols-[1fr_1fr_150px_100px_80px] gap-4 px-4 py-3 bg-muted/50 border-b border-border text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                <span>Name</span>
                <span>Contact</span>
                <span>Role</span>
                <span>Status</span>
                <span></span>
              </div>

              {/* Table Rows */}
              {users.map((user) => (
                <div
                  key={user.id}
                  className="grid grid-cols-[1fr_1fr_150px_100px_80px] gap-4 px-4 py-3 border-b border-border last:border-b-0 hover:bg-muted/30 transition-colors items-center"
                >
                  <div>
                    <p className="text-sm font-medium text-foreground">{user.name}</p>
                  </div>
                  <div>
                    <p className="text-sm text-foreground">{user.email}</p>
                    <p className="text-xs text-muted-foreground">{user.phone}</p>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <Shield className="w-3.5 h-3.5 text-accent" />
                    <span className="text-sm text-foreground">{user.role}</span>
                  </div>
                  <div>
                    <span
                      className={`inline-flex items-center px-2 py-0.5 text-xs font-medium rounded-full ${
                        user.status === "active"
                          ? "bg-accent/20 text-accent"
                          : "bg-muted text-muted-foreground"
                      }`}
                    >
                      {user.status}
                    </span>
                  </div>
                  <div className="flex items-center gap-1 justify-end">
                    <button
                      onClick={() => handleEdit(user)}
                      className="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
                      aria-label={`Edit ${user.name}`}
                    >
                      <Pencil className="w-3.5 h-3.5" />
                    </button>
                    <button
                      className="p-1.5 rounded-md text-destructive/70 hover:text-destructive hover:bg-destructive/10 transition-colors"
                      aria-label={`Delete ${user.name}`}
                    >
                      <Trash2 className="w-3.5 h-3.5" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Create User Modal */}
      <CreateUserModal
        open={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        onSuccess={() => setShowCreateModal(false)}
      />

      {/* Edit User Modal */}
      <EditUserModal
        open={showEditModal}
        onClose={() => {
          setShowEditModal(false);
          setEditingUser(null);
        }}
        onSuccess={() => {
          setShowEditModal(false);
          setEditingUser(null);
        }}
        user={editingUser}
      />
    </>
  );
}
