export const validatePhone = (value) => {
    if (!value) return '';
    if (/\s/.test(value)) return 'Số điện thoại không được chứa khoảng trắng';
    if (/[a-zA-Z]/.test(value)) return 'Số điện thoại không được chứa chữ cái';
    if (/[^0-9]/.test(value)) return 'Số điện thoại không được chứa ký tự đặc biệt';
    if (!value.startsWith('0')) return 'Số điện thoại phải bắt đầu bằng số 0';
    if (value.length > 10) return 'Số điện thoại không được quá 10 số';
    if (value.length < 10) return 'Số điện thoại phải đủ 10 số';
    return '';
};

export const handlePhoneChange = (rawValue) => {
    const digitsOnly = rawValue.replace(/\D/g, '').slice(0, 10);
    let error = '';
    if (/\s/.test(rawValue)) error = 'Số điện thoại không được chứa khoảng trắng';
    else if (/[a-zA-Z]/.test(rawValue)) error = 'Số điện thoại không được chứa chữ cái';
    else if (/[^0-9]/.test(rawValue)) error = 'Số điện thoại không được chứa ký tự đặc biệt';
    else if (digitsOnly && !digitsOnly.startsWith('0')) error = 'Số điện thoại phải bắt đầu bằng số 0';
    else if (rawValue.replace(/\D/g, '').length > 10) error = 'Số điện thoại không được quá 10 số';
    return { cleaned: digitsOnly, error };
};

export const validateEmail = (email) => {
    if (!email) return '';
    if (/\s/.test(email)) return 'Email không được chứa khoảng trắng';
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) return 'Định dạng email không hợp lệ (ví dụ: abc@gmail.com)';
    return '';
};

export const handleEmailChange = (rawValue) => {
    const cleaned = rawValue.replace(/\s/g, '');
    let error = '';
    if (/\s/.test(rawValue)) error = 'Email không được chứa khoảng trắng';
    else error = validateEmail(cleaned);
    return { cleaned, error };
};

export const validateName = (name) => {
    if (!name) return '';
    const nameRegex = /^[a-zA-ZÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚĂĐĨŨƠàáâãèéêìíòóôõùúăđĩũơƯĂẠẢẤẦẨẪẬẮẰẲẴẶẸẺẼỀỀỂẾưăạảấầẩẫậắằẳẵặẹẻẽềềểếỄỆỈỊỌỎỐỒỔỖỘỚỜỞỠỢỤỦỨỪễệỉịọỏốồổỗộớờởỡợụủứừỬỮỰỲỴÝỶỸửữựỳỵỷỹ\s]+$/;
    if (!nameRegex.test(name)) return 'Họ tên không được chứa số hoặc ký tự đặc biệt';
    return '';
};

export const handleNameChange = (rawValue) => {
    let error = '';
    const nameRegex = /^[a-zA-ZÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚĂĐĨŨƠàáâãèéêìíòóôõùúăđĩũơƯĂẠẢẤẦẨẪẬẮẰẲẴẶẸẺẼỀỀỂẾưăạảấầẩẫậắằẳẵặẹẻẽềềểếỄỆỈỊỌỎỐỒỔỖỘỚỜỞỠỢỤỦỨỪễệỉịọỏốồổỗộớờởỡợụủứừỬỮỰỲỴÝỶỸửữựỳỵỷỹ\s]+$/;
    if (rawValue && !nameRegex.test(rawValue)) {
        error = 'Họ tên không được chứa số hoặc ký tự đặc biệt';
    }
    return { cleaned: rawValue, error };
};
